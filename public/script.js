let urlDatabase = {};
let urlCounter = 0;
let totalClicks = 0;

function generateShortCode() {
  const chars =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  let result = "";
  for (let i = 0; i < 6; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

function isValidUrl(string) {
  try {
    new URL(string);
    return true;
  } catch (_) {
    return false;
  }
}

function shortenUrl() {
  const urlInput = document.getElementById("urlInput");
  const longUrl = urlInput.value.trim();

  if (!longUrl) {
    alert("Please enter a URL!");
    return;
  }

  if (!isValidUrl(longUrl)) {
    alert("Please enter a valid URL (include http:// or https://)");
    return;
  }

  const shortCode = generateShortCode();
  const shortUrl = `https://short.ly/${shortCode}`;

  // Store in our "database"
  urlDatabase[shortCode] = {
    originalUrl: longUrl,
    clicks: 0,
    created: new Date(),
  };

  // Update stats
  urlCounter++;
  document.getElementById("urlCount").textContent = urlCounter;

  // Show result
  document.getElementById("shortUrl").textContent = shortUrl;
  document.getElementById("result").classList.add("show");

  // Clear input
  urlInput.value = "";
}

function copyToClipboard() {
  const shortUrl = document.getElementById("shortUrl").textContent;
  const copyBtn = document.getElementById("copyBtn");

  navigator.clipboard
    .writeText(shortUrl)
    .then(() => {
      copyBtn.textContent = "✅ Copied!";
      copyBtn.classList.add("copied");

      // Simulate a click (in real app, this would happen when someone visits the short URL)
      totalClicks++;
      document.getElementById("clickCount").textContent = totalClicks;

      setTimeout(() => {
        copyBtn.textContent = "📋 Copy Link";
        copyBtn.classList.remove("copied");
      }, 2000);
    })
    .catch(() => {
      // Fallback for older browsers
      const textArea = document.createElement("textarea");
      textArea.value = shortUrl;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand("copy");
      document.body.removeChild(textArea);

      copyBtn.textContent = "✅ Copied!";
      copyBtn.classList.add("copied");

      setTimeout(() => {
        copyBtn.textContent = "📋 Copy Link";
        copyBtn.classList.remove("copied");
      }, 2000);
    });
}

// Allow Enter key to trigger shortening
document.getElementById("urlInput").addEventListener("keypress", function (e) {
  if (e.key === "Enter") {
    shortenUrl();
  }
});
