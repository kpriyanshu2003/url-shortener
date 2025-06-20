async function shortenUrl() {
  const url = document.getElementById("urlInput").value.trim();
  if (!url) return showToast("Please enter a URL!");
  if (!isValidUrl(url))
    return showToast("Please enter a valid URL (include http:// or https://)");

  try {
    const res = await fetch("/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url }),
    });

    const data = await res.json();
    if (!data?.short_url)
      return showToast("Failed to shorten URL. Please try again.");

    const short = `${window.location.origin}${data.short_url}`;
    document.getElementById("shortUrl").textContent = short;
    document.getElementById("result").classList.add("show");
    document.getElementById("urlInput").value = "";
  } catch (err) {
    console.error("Error:", err);
    showToast("Error connecting to the server.");
  }
}

function copyToClipboard() {
  const text = document.getElementById("shortUrl").textContent;
  navigator.clipboard
    .writeText(text)
    .then(() => {
      document.getElementById("clickCount").textContent = ++totalClicks;
      showToast("✅ Copied!");
    })
    .catch(() => {
      const t = document.createElement("textarea");
      t.value = text;
      document.body.appendChild(t);
      t.select();
      document.execCommand("copy");
      document.body.removeChild(t);
      showToast("✅ Copied!");
    });
}

const isValidUrl = (url) => {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
};

function showToast(msg) {
  if (!msg) return;
  const x = document.getElementById("snackbar");
  x.textContent = msg;
  x.className = "show";
  setTimeout(() => x.classList.remove("show"), 3000);
}

document.getElementById("urlInput").addEventListener("keypress", (e) => {
  if (e.key === "Enter") shortenUrl();
});

document.addEventListener("keydown", (e) => {
  if (e.ctrlKey && e.key === "c") {
    copyToClipboard();
    e.preventDefault();
    showToast("✅ Copied!");
  }
});
