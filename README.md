# Mathilde

Adaptive math learning app that grows with the learner.

## Mission

Build math confidence through short, focused practice sessions that adapt to the learner's understanding. No punishment, no streaks, no shame — just steady progress at their pace.

## How it works

**Two personas:**
- **Curator** — adds topics, reviews AI-generated sessions, monitors progress. All via CLI.
- **Learner** — opens the app, works through a 10-15 minute session, earns points and levels. Always moves forward.

**The loop:**
1. Curator adds a topic → AI breaks it into concepts with prerequisites
2. Planner generates a session (complete HTML with exercises, hints, feedback) → Curator reviews and approves
3. Learner completes the session → results flow to Firestore
4. AI analyzes results, updates learning records with what the learner knows and where they struggle
5. Planner uses updated records to generate the next session at the right difficulty
6. Repeat

## Architecture

```
GitHub Pages (static SPA, vanilla JS)
├── Auth (Firebase, Google sign-in)
├── Renders AI-generated HTML sessions
└── Bridge API: window.mathilde.*

Firebase (Firestore)
├── Sessions, results, learning records
├── Profile (XP, level, mission)
└── Topics with concept trees

Go backend (home machine)
├── Reactor — watches Firestore, runs post-session analysis
├── Planner — generates sessions via AI
└── CLI — curator interface (mathilde topic/session/queue/progress)
```

## Learning philosophy

The app is built around the idea that learning works best when it respects how attention and memory actually work:

- **Short sessions** — 10-15 minutes of focused practice beats an hour of drifting attention. The learner can keep going if they're in the zone.
- **Forward-only progress** — points and levels only go up. A bad day doesn't erase progress. Showing up and trying always counts.
- **Scaffolding over failure** — every exercise has a hint chain that progressively simplifies the problem. The learner always arrives at understanding — sometimes with more support, sometimes less.
- **Low cognitive load** — progress bars show remaining effort, layouts are clean, feedback is immediate. No text walls, no clutter.
- **Adaptive difficulty** — frustration signals (many hints, fast wrong answers) feed into the next session's planning, not a mid-session pivot. The system adjusts without the learner noticing.

## Pedagogy

Built on Matt Pocock's [teach skill](https://github.com/mattpocock/skills) philosophy:
- Mission-grounded learning tied to real goals
- Zone of Proximal Development — always challenged just enough
- Storage strength over fluency — spaced repetition, retrieval practice
- Knowledge taught first, then practiced as skills
- Tight feedback loops with specific, helpful feedback

AI drives all pedagogy. The backend is plumbing.
