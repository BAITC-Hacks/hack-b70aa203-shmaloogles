{
  description = "MOST development environment: Go, Nuxt and PostgreSQL";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in {
      devShells = forAllSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gopls
              nodejs_22
              postgresql_16
              git
              curl
              jq
              docker-client
              docker-compose
            ];

            # Defaults only: preserve explicitly supplied settings and secrets.
            # Enter from the repository root with: nix develop
            # The Docker daemon must be installed/running on the host.
            # This shell does not start services or execute .env.
            shellHook = ''
              export AI_MODE="''${AI_MODE:-mock}"
              export DATABASE_URL="''${DATABASE_URL:-postgres://shmaloogles:shmaloogles@localhost:5432/shmaloogles?sslmode=disable}"
              export CORS_ORIGIN="''${CORS_ORIGIN:-http://localhost:3000}"
              export NUXT_API_BASE="''${NUXT_API_BASE:-http://localhost:8080}"
              export NUXT_PUBLIC_API_BASE="''${NUXT_PUBLIC_API_BASE:-http://localhost:8080}"

              printf '%s\n' \
                'MOST development shell' \
                'Database: docker compose up -d postgres' \
                'Backend:  cd backend && go run ./cmd/api' \
                'Frontend: cd frontend && npm ci && npm run dev' \
                'Run backend and frontend in separate terminals.' \
                'Commit the generated flake.lock to pin the toolchain.'
            '';
          };
        });
    };
}
