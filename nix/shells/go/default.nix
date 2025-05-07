{ mkShell, pkgs, ... }:
mkShell {
  packages = with pkgs; [
    go
    pre-commit
    detect-secrets
    gnused
    sshpass
    go-swag
  ];
}
