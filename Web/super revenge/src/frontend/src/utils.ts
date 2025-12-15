export const isPrivateIp = (raw: string | undefined | null): boolean => {
  if (!raw) return false;

  let ip = raw;

  // Strip IPv4-mapped IPv6 prefix, e.g. ::ffff:192.168.1.5
  if (ip.startsWith('::ffff:')) {
    ip = ip.substring('::ffff:'.length);
  }

  // Strip port if present in IPv4 "host:port" style
  const ipv4PortMatch = ip.match(/^(\d+\.\d+\.\d+\.\d+):\d+$/);
  if (ipv4PortMatch) {
    ip = ipv4PortMatch[1];
  }

  // Loopback
  if (ip === '127.0.0.1' || ip === '::1') {
    return true;
  }

  // IPv4 private ranges (RFC1918)
  const parts = ip.split('.');
  if (parts.length === 4) {
    const [a, b, c, d] = parts.map((x) => Number(x));
    if ([a, b, c, d].some((x) => Number.isNaN(x) || x < 0 || x > 255)) {
      return false;
    }

    // 10.0.0.0/8
    if (a === 10) return true;

    // 172.16.0.0/12  (172.16.0.0 – 172.31.255.255)
    if (a === 172 && b >= 16 && b <= 31) return true;

    // 192.168.0.0/16
    if (a === 192 && b === 168) return true;

    return false;
  }

  // Very basic IPv6 "local enough" check
  // ULA: fc00::/7 (fc, fd)
  // Link-local: fe80::/10
  const lower = ip.toLowerCase();
  if (
    lower.startsWith('fc') ||
    lower.startsWith('fd') ||
    lower.startsWith('fe80:')
  ) {
    return true;
  }

  return false;
};


export function genHex() {
  const n = Math.floor(Math.random() * 2 ** 32);
  return n.toString(16).padStart(8, "0");
}

export function genBase36() {
  const n = Math.floor(Math.random() * 36 ** 8);
  return n.toString(36).padStart(8, "0");
}

export function genBase64() {
  const array = new Uint8Array(6);
  crypto.getRandomValues(array);
  return btoa(String.fromCharCode(...array)).replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_');
}

