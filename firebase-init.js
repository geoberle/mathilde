(function () {
  "use strict";

  var config = {
    apiKey: "AIzaSyBehCsEWPjiTI-M-s35yTLBSFtN97mKb5o",
    authDomain: "mathilde-61d77.firebaseapp.com",
    projectId: "mathilde-61d77",
    storageBucket: "mathilde-61d77.firebasestorage.app",
    messagingSenderId: "457178130197",
    appId: "1:457178130197:web:29c3d76fdacf538a0fab1c"
  };

  firebase.initializeApp(config);

  var auth = firebase.auth();
  var firestore = firebase.firestore();

  firestore.enablePersistence().catch(function (err) {
    if (err.code === "failed-precondition") {
      console.warn("Firestore persistence unavailable: multiple tabs open");
    } else if (err.code === "unimplemented") {
      console.warn("Firestore persistence unavailable: browser not supported");
    }
  });

  function userDoc() {
    var user = auth.currentUser;
    if (!user) return null;
    return firestore.collection("users").doc(user.uid);
  }

  function profileDoc() {
    var doc = userDoc();
    if (!doc) return null;
    return doc.collection("profile").doc("main");
  }

  function xpForLevel(level) {
    return 100 + (level * 20);
  }

  window.mathilde = {
    auth: auth,
    firestore: firestore,

    signInWithGoogle: function () {
      var provider = new firebase.auth.GoogleAuthProvider();
      return auth.signInWithPopup(provider);
    },

    signOut: function () {
      return auth.signOut();
    },

    onAuthStateChanged: function (cb) {
      return auth.onAuthStateChanged(cb);
    },

    getProfile: function () {
      var doc = profileDoc();
      if (!doc) return Promise.resolve({ xp: 0, level: 1 });
      return doc.get().then(function (snap) {
        if (!snap.exists) return { xp: 0, level: 1 };
        var d = snap.data();
        return { xp: d.xp || 0, level: d.level || 1 };
      }).catch(function () { return { xp: 0, level: 1 }; });
    },

    addXP: function (amount) {
      var doc = profileDoc();
      if (!doc || amount <= 0) return Promise.resolve(null);
      return firestore.runTransaction(function (tx) {
        return tx.get(doc).then(function (snap) {
          var data = snap.exists ? snap.data() : {};
          var currentXP = data.xp || 0;
          var currentLevel = data.level || 1;
          var newXP = currentXP + amount;
          var newLevel = currentLevel;
          while (newXP >= xpForLevel(newLevel)) {
            newXP -= xpForLevel(newLevel);
            newLevel++;
          }
          tx.set(doc, { xp: newXP, level: newLevel }, { merge: true });
          return { oldLevel: currentLevel, newLevel: newLevel, xp: newXP, leveledUp: newLevel > currentLevel };
        });
      });
    },

    xpForLevel: xpForLevel,

    submitAnswer: function (exerciseId, result) {
      // Bridge API stub — sessions will call this
    },

    sessionComplete: function (summary) {
      // Bridge API stub — sessions will call this
    },

    requestEvaluation: function (imageData) {
      // Bridge API stub — handwriting evaluation (post-MVP)
      return Promise.resolve(null);
    }
  };
})();
