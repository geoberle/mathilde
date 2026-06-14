(function () {
  "use strict";

  var app = document.getElementById("app");
  var userInfo = document.getElementById("user-info");

  function esc(str) {
    var div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
  }

  function renderSignIn() {
    userInfo.innerHTML = "";
    app.innerHTML =
      '<div class="welcome-anon">' +
        '<p class="welcome-tagline">Deine Mathe-Lernapp</p>' +
        '<button class="sign-in-btn" id="sign-in-btn">Mit Google anmelden</button>' +
      '</div>';
    document.getElementById("sign-in-btn").addEventListener("click", function () {
      window.mathilde.signInWithGoogle().catch(function (err) {
        if (err.code !== "auth/popup-closed-by-user") {
          console.error("Sign-in error:", err);
        }
      });
    });
  }

  function renderWelcome(user) {
    var displayName = user.displayName || "Lernende";

    userInfo.innerHTML =
      '<span class="user-name">' + esc(displayName) + '</span>' +
      '<button class="sign-out-btn" id="sign-out-btn">Abmelden</button>';

    app.innerHTML =
      '<div class="welcome">' +
        '<h2>Hallo, ' + esc(displayName.split(" ")[0]) + '!</h2>' +
        '<p class="welcome-sub">Keine Sessions verfügbar. Dein Curator bereitet neue Inhalte vor.</p>' +
        '<div class="stats" id="stats"></div>' +
      '</div>';

    document.getElementById("sign-out-btn").addEventListener("click", function () {
      window.mathilde.signOut();
    });

    window.mathilde.getProfile().then(function (profile) {
      var statsEl = document.getElementById("stats");
      if (!statsEl) return;
      statsEl.innerHTML =
        '<div class="stat">' +
          '<span class="stat-label">Level</span>' +
          '<span class="stat-value">' + profile.level + '</span>' +
        '</div>' +
        '<div class="stat">' +
          '<span class="stat-label">XP</span>' +
          '<span class="stat-value">' + profile.xp + ' / ' + window.mathilde.xpForLevel(profile.level) + '</span>' +
        '</div>';
    });
  }

  function renderLoading() {
    app.innerHTML = '<div class="loading">Laden…</div>';
  }

  renderLoading();

  window.mathilde.onAuthStateChanged(function (user) {
    if (user) {
      renderWelcome(user);
    } else {
      renderSignIn();
    }
  });
})();
