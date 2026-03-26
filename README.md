# 📦 Pack Planner

This repository is organized as two sibling applications:

- [`packplanner/`](/Users/melih/work/case/repartners/packplanner)
  Go backend API, business logic, Swagger docs, and Docker setup
- [`ui/`](/Users/melih/work/case/repartners/ui)
  Standalone frontend that can be deployed separately, for example on Amplify

The challenge rules are implemented in this order:

1. Packs cannot be split.
2. The shipped total must be as small as possible while still fulfilling the order.
3. If multiple solutions ship the same total, the one with fewer packs wins.

Example:

- Order `501`
- Pack sizes `[250, 500, 1000, 2000, 5000]`
- Best result: `500 + 250 = 750`

`1000` would use fewer packs, but it ships more items than necessary, so it is not the correct answer.

## Live Demo

- Frontend: [https://packplanner-ui.onrender.com](https://packplanner-ui.onrender.com)
- Swagger UI: [https://packplanner.onrender.com/swagger](https://packplanner.onrender.com/swagger)
- Demo video: [Watch the walkthrough](https://github.com/user-attachments/assets/4b6f332b-5e6a-4eac-9365-2d1b5e23463c)

## Repository Layout

```text
repo-root/
├── packplanner/
│   ├── cmd/api
│   ├── internal
│   ├── docs
│   ├── Dockerfile
│   └── README.md
└── ui/
    ├── app.js
    ├── config.js
    ├── index.html
    └── styles.css
```

## Architecture

<img src="packplanner/docs/infra-diagram.png" alt="PackPlanner infrastructure diagram" width="720" />

High-level flow:

- The frontend calls the backend HTTP API.
- Echo routes requests to handlers.
- The application layer coordinates use cases.
- The domain planner calculates the optimal pack mix.
- The repository provides the active pack size configuration.

## Backend

Backend source lives in:

- [`packplanner/`](/Users/melih/work/case/repartners/packplanner)

Run locally:

```bash
cd packplanner
go run ./cmd/api
```

Run tests:

```bash
cd packplanner
go test ./...
```

API base URL:

- [http://localhost:8680](http://localhost:8680)
- Swagger: [http://localhost:8680/swagger](http://localhost:8680/swagger)

## Frontend

Frontend source lives in:

- [`ui/`](/Users/melih/work/case/repartners/ui)

The frontend is intentionally independent from the Go service:

- It is not served by the backend.
- It can be deployed on Amplify or any static host.
- It only needs the backend API URL in `ui/config.js`.

Example config:

```js
window.PACKPLANNER_CONFIG = {
  apiBaseUrl: "http://localhost:8680",
};
```

Run locally with a static server:

```bash
npx serve ./ui -l 3000
```

Then open:

- [http://localhost:3000](http://localhost:3000)

## Deployment Idea

- Deploy `packplanner/` as the backend service
- Deploy `ui/` as a static frontend
- Point `ui/config.js` to the backend domain

## More Details

For backend-specific implementation notes, API examples, Docker usage, and algorithm details:

- [packplanner/README.md](/Users/melih/work/case/repartners/packplanner/README.md)
