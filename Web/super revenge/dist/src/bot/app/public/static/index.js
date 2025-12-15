window.turnstileToken = "";
window.onTurnstileSuccess = (token) => {
  window.turnstileToken = token;
};
window.onTurnstileExpired = () => {
  window.turnstileToken = "";
};

const clickSfx = document.getElementById("click-sfx");

window.addEventListener("DOMContentLoaded", () => {
  const realInput = document.getElementById("urlInput");
  const fakeInput = document.getElementById("fakeInput");

  fakeInput.readOnly = true;

  document.addEventListener("keydown", (e) => {
    if (e.ctrlKey || e.altKey || e.metaKey) return;

    if (e.key === "Backspace") {
      realInput.value = realInput.value.slice(0, -1);
    } else if (e.key === "Enter") {
      startCrawler();
    } else if (e.key.length === 1 && /[\x20-\x7E]/.test(e.key)) {
      realInput.value += e.key;
    }

    fakeInput.value = toKatakanaDisplay(realInput.value);

    e.preventDefault();
  });

  fakeInput.addEventListener("paste", (e) => {
    const paste = (e.clipboardData || window.clipboardData).getData("text");
    const clean = paste.replace(/[^\x20-\x7E]/g, ""); // remove non-ASCII
    realInput.value += clean;
    fakeInput.value = toKatakanaDisplay(realInput.value);
    e.preventDefault();
  });

  function toKatakanaDisplay(text) {
    const map = {
      a: "ア",
      b: "ビ",
      c: "ク",
      d: "デ",
      e: "エ",
      f: "フ",
      g: "グ",
      h: "ホ",
      i: "イ",
      j: "ジ",
      k: "カ",
      l: "ル",
      m: "ム",
      n: "ヌ",
      o: "オ",
      p: "プ",
      q: "ク",
      r: "ラ",
      s: "ス",
      t: "ト",
      u: "ウ",
      v: "ヴ",
      w: "ワ",
      x: "クス",
      y: "ヨ",
      z: "ゼ",
      A: "ア",
      B: "ビ",
      C: "ク",
      D: "デ",
      E: "エ",
      F: "フ",
      G: "グ",
      H: "ホ",
      I: "イ",
      J: "ジ",
      K: "カ",
      L: "ル",
      M: "ム",
      N: "ヌ",
      O: "オ",
      P: "プ",
      Q: "ク",
      R: "ラ",
      S: "ス",
      T: "ト",
      U: "ウ",
      V: "ヴ",
      W: "ワ",
      X: "クス",
      Y: "ヨ",
      Z: "ゼ",
      " ": "　",
      "/": "／",
      ".": "。",
      ":": "：",
      "-": "ー",
      "?": "？",
      "=": "＝",
      "&": "＆",
    };

    return [...text].map((char) => map[char] || char).join("");
  }

  const container = document.getElementById("personaTitle");
  putText("Ubermensch");

  function putText(text) {
    const emojis = [
      "⚔",
      "愛",
      "忍",
      "侍",
      "魂",
      "影",
      "龍",
      "刃",
      "鬼",
      "怒",
      "花",
      "雪",
      "月",
    ];

    const line = document.createElement("div");
    line.style.display = "flex";
    line.style.justifyContent = "center";
    line.style.marginBottom = "1rem";
    line.style.flexWrap = "wrap";
    line.style.gap = "0.5rem";

    for (let i = 0; i < text.length; i++) {
      const emoji = emojis[i % emojis.length];
      const char = text[i];

      const box = document.createElement("div");
      box.className = "emoji-char-box";
      box.textContent = `${emoji}\n${char}`;
      box.style.whiteSpace = "pre";
      box.style.display = "inline-block";
      box.style.textAlign = "center";
      box.style.margin = "0 6px";
      box.style.border = "2.5px solid #ff0000";
      box.style.background = "#111";
      box.style.color = "#fff";
      box.style.padding = "10px 12px";
      box.style.fontFamily = "Impact, Arial Black, monospace";
      box.style.fontSize = "1.1rem";
      box.style.lineHeight = "1.4";
      box.style.transform = `rotate(${rand(-5, 5)}deg) skew(${rand(
        -6,
        6
      )}deg, ${rand(-6, 6)}deg) translate(${rand(-2, 2)}px, ${rand(-2, 2)}px)`;

      line.appendChild(box);
    }

    container.appendChild(line);
  }

  function rand(min, max) {
    return Math.floor(Math.random() * (max  - min + 1)) + min;
  }

  const visitBtn = document.getElementById("visitBtn");
  visitBtn.addEventListener("click", startCrawler);
});

async function startCrawler() {
  if (
    document.querySelector(".cf-turnstile") &&
    !window.turnstileToken
  ) {
    alert("Complete the captcha first.");
    return;
  }

  document.body.classList.remove("persona-bg");
  document.body.classList.add("persona-bg-alt");

  clickSfx.currentTime = 0;
  clickSfx.play();

  const modal = document.getElementById("statusModal");
  modal.style.display = "block";

  const modalStatus = document.getElementById("modalStatus");
  const progressBar = document.getElementById("modalProgressBar");

  modalStatus.innerText = "Sending request to server...";

  const csrfToken = document.getElementById("csrfToken").value;

  progressBar.style.width = "40%";

  const res = await fetch(window.location.pathname, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "CSRF-Token": csrfToken,
    },
    body: JSON.stringify({
      action: "run",
      url: document.getElementById("urlInput").value,
      turnstileToken: window.turnstileToken,
    }),
  });

  modalStatus.innerText = "Processing response...";
  progressBar.style.width = "70%";

  const json = await res.json();

  modalStatus.innerText = json.message;
  progressBar.style.width = "100%";

  status.innerText = json.message;

  if (window.turnstile?.reset) {
    window.turnstile.reset();
  }
  window.turnstileToken = "";

  if (json.message === "Crawl complete.") {
    document.body.classList.remove("persona-bg-alt");
    document.body.classList.add("persona-bg");
    setTimeout(() => {
      modal.style.display = "none";
      progressBar.style.width = "25%";
      modalStatus.innerText = "Initializing...";
    }, 2000);
  }
}
