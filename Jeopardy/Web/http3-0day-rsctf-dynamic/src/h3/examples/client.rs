use std::{net::IpAddr, path::PathBuf, sync::Arc};

use futures::future;
use h3::error::{ConnectionError, StreamError};
use rustls::{
    client::danger::{HandshakeSignatureValid, ServerCertVerified, ServerCertVerifier},
    pki_types::{CertificateDer, ServerName, UnixTime},
    CertificateError,
    DigitallySignedStruct,
    SignatureScheme,
};
use rustls_native_certs::CertificateResult;
use structopt::StructOpt;
use tokio::io::AsyncWriteExt;
use tracing::{error, info};

use h3_quinn::quinn;

static ALPN: &[u8] = b"h3";

#[derive(StructOpt, Debug)]
#[structopt(name = "server")]
struct Opt {
    #[structopt(
        long,
        short,
        default_value = "examples/server.cert",
        help = "Certificate of CA who issues the server certificate"
    )]
    pub ca: PathBuf,

    #[structopt(name = "keylogfile", long)]
    pub key_log_file: bool,

    #[structopt(long, help = "Disable TLS certificate verification (dev only)")]
    pub insecure: bool,

    #[structopt()]
    pub uri: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt()
        .with_env_filter(tracing_subscriber::EnvFilter::from_default_env())
        .with_span_events(tracing_subscriber::fmt::format::FmtSpan::FULL)
        .with_writer(std::io::stderr)
        .with_max_level(tracing::Level::INFO)
        .init();

    let opt = Opt::from_args();

    // DNS lookup

    let uri = opt.uri.parse::<http::Uri>()?;

    if uri.scheme() != Some(&http::uri::Scheme::HTTPS) {
        Err("uri scheme must be 'https'")?;
    }

    let auth = uri.authority().ok_or("uri must have a host")?.clone();

    let port = auth.port_u16().unwrap_or(443);

    let addr = tokio::net::lookup_host((auth.host(), port))
        .await?
        .next()
        .ok_or("dns found no addresses")?;

    info!("DNS lookup for {:?}: {:?}", uri, addr);

    // create quinn client endpoint

    // load CA certificates stored in the system
    let mut roots = rustls::RootCertStore::empty();
    let CertificateResult { certs, errors, .. } = rustls_native_certs::load_native_certs();
    for cert in certs {
        if let Err(e) = roots.add(cert) {
            error!("failed to parse trust anchor: {}", e);
        }
    }
    for e in errors {
        error!("couldn't load default trust roots: {}", e);
    }

    // load certificate of CA who issues the server certificate
    // NOTE that this should be used for dev only
    let ca_bytes = std::fs::read(&opt.ca)?;
    let mut pinned_verifier: Option<Arc<PinnedCertVerifier>> = None;
    if let Err(e) = roots.add(CertificateDer::from(ca_bytes.clone())) {
        error!("failed to parse trust anchor (will pin instead): {}", e);
        pinned_verifier = Some(Arc::new(PinnedCertVerifier {
            expected: CertificateDer::from(ca_bytes),
        }));
    }

    let mut tls_config = rustls::ClientConfig::builder()
        .with_root_certificates(roots)
        .with_no_client_auth();

    if opt.insecure {
        tls_config
            .dangerous()
            .set_certificate_verifier(Arc::new(NoVerifier));
    } else if let Some(pinned) = pinned_verifier {
        tls_config.dangerous().set_certificate_verifier(pinned);
    }

    tls_config.enable_early_data = true;
    tls_config.alpn_protocols = vec![ALPN.into()];

    // optional debugging support
    if opt.key_log_file {
        // Write all Keys to a file if SSLKEYLOGFILE is set
        // WARNING, we enable this for the example, you should think carefully about enabling in your own code
        tls_config.key_log = Arc::new(rustls::KeyLogFile::new());
    }

    let mut client_endpoint = h3_quinn::quinn::Endpoint::client("[::]:0".parse().unwrap())?;

    let client_config = quinn::ClientConfig::new(Arc::new(
        quinn::crypto::rustls::QuicClientConfig::try_from(tls_config)?,
    ));
    client_endpoint.set_default_client_config(client_config);

    let server_name_for_sni = match auth.host().parse::<IpAddr>() {
        Ok(_) => "localhost",
        Err(_) => auth.host(),
    };

    let conn = client_endpoint.connect(addr, server_name_for_sni)?.await?;

    info!("QUIC connection established");

    // create h3 client

    // h3 is designed to work with different QUIC implementations via
    // a generic interface, that is, the [`quic::Connection`] trait.
    // h3_quinn implements the trait w/ quinn to make it work with h3.
    let quinn_conn = h3_quinn::Connection::new(conn);

    let (mut driver, mut send_request) = h3::client::new(quinn_conn).await?;

    let drive = async move {
        return Err::<(), ConnectionError>(future::poll_fn(|cx| driver.poll_close(cx)).await);
    };

    // In the following block, we want to take ownership of `send_request`:
    // the connection will be closed only when all `SendRequest`s instances
    // are dropped.
    //
    //             So we "move" it.
    //                  vvvv
    let request = async move {
        info!("sending request ...");

        let req = http::Request::builder().uri(uri).body(()).unwrap();

        // sending request results in a bidirectional stream,
        // which is also used for receiving response
        let mut stream = send_request.send_request(req).await?;

        // finish on the sending side
        stream.finish().await?;

        info!("receiving response ...");

        let resp = stream.recv_response().await?;

        info!("response: {:?} {}", resp.version(), resp.status());
        info!("headers: {:#?}", resp.headers());

        // `recv_data()` must be called after `recv_response()` for
        // receiving potential response body
        while let Some(mut chunk) = stream.recv_data().await? {
            let mut out = tokio::io::stdout();
            out.write_all_buf(&mut chunk).await.unwrap();
            out.flush().await.unwrap();
        }

        Ok::<_, StreamError>(())
    };

    let (req_res, drive_res) = tokio::join!(request, drive);

    if let Err(err) = req_res {
        if err.is_h3_no_error() {
            info!("connection closed with H3_NO_ERROR");
        } else {
            error!("request failed: {:?}", err);
        }
        error!("request failed: {:?}", err);
    }
    if let Err(err) = drive_res {
        if err.is_h3_no_error() {
            info!("connection closed with H3_NO_ERROR");
        } else {
            error!("connection closed with error: {:?}", err);
            return Err(err.into());
        }
    }

    // wait for the connection to be closed before exiting
    client_endpoint.wait_idle().await;

    Ok(())
}

#[derive(Debug)]
struct NoVerifier;

impl ServerCertVerifier for NoVerifier {
    fn verify_server_cert(
        &self,
        _end_entity: &CertificateDer<'_>,
        _intermediates: &[CertificateDer<'_>],
        _server_name: &ServerName<'_>,
        _ocsp_response: &[u8],
        _now: UnixTime,
    ) -> Result<ServerCertVerified, rustls::Error> {
        Ok(ServerCertVerified::assertion())
    }

    fn verify_tls12_signature(
        &self,
        _message: &[u8],
        _cert: &CertificateDer<'_>,
        _dss: &DigitallySignedStruct,
    ) -> Result<HandshakeSignatureValid, rustls::Error> {
        Ok(HandshakeSignatureValid::assertion())
    }

    fn verify_tls13_signature(
        &self,
        _message: &[u8],
        _cert: &CertificateDer<'_>,
        _dss: &DigitallySignedStruct,
    ) -> Result<HandshakeSignatureValid, rustls::Error> {
        Ok(HandshakeSignatureValid::assertion())
    }

    fn supported_verify_schemes(&self) -> Vec<SignatureScheme> {
        vec![
            SignatureScheme::ECDSA_NISTP256_SHA256,
            SignatureScheme::ECDSA_NISTP384_SHA384,
            SignatureScheme::ED25519,
            SignatureScheme::RSA_PSS_SHA256,
            SignatureScheme::RSA_PKCS1_SHA256,
        ]
    }
}

#[derive(Debug)]
struct PinnedCertVerifier {
    expected: CertificateDer<'static>,
}

impl ServerCertVerifier for PinnedCertVerifier {
    fn verify_server_cert(
        &self,
        end_entity: &CertificateDer<'_>,
        _intermediates: &[CertificateDer<'_>],
        _server_name: &ServerName<'_>,
        _ocsp_response: &[u8],
        _now: UnixTime,
    ) -> Result<ServerCertVerified, rustls::Error> {
        if end_entity == &self.expected {
            Ok(ServerCertVerified::assertion())
        } else {
            Err(rustls::Error::InvalidCertificate(
                CertificateError::ApplicationVerificationFailure,
            ))
        }
    }

    fn verify_tls12_signature(
        &self,
        _message: &[u8],
        _cert: &CertificateDer<'_>,
        _dss: &DigitallySignedStruct,
    ) -> Result<HandshakeSignatureValid, rustls::Error> {
        Ok(HandshakeSignatureValid::assertion())
    }

    fn verify_tls13_signature(
        &self,
        _message: &[u8],
        _cert: &CertificateDer<'_>,
        _dss: &DigitallySignedStruct,
    ) -> Result<HandshakeSignatureValid, rustls::Error> {
        Ok(HandshakeSignatureValid::assertion())
    }

    fn supported_verify_schemes(&self) -> Vec<SignatureScheme> {
        vec![
            SignatureScheme::ECDSA_NISTP256_SHA256,
            SignatureScheme::ECDSA_NISTP384_SHA384,
            SignatureScheme::ED25519,
            SignatureScheme::RSA_PSS_SHA256,
            SignatureScheme::RSA_PKCS1_SHA256,
        ]
    }
}
