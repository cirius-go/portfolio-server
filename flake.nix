{
  description = "portfolio-server Backend Flake";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    darwin = {
      url = "github:lnl7/nix-darwin";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    snowfall-lib = {
      url = "github:snowfallorg/lib";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs:
    let
      lib = inputs.snowfall-lib.mkLib {
        inherit inputs;
        src = ./.;
        snowfall = {
          root = ./nix;
          meta = {
            name = "portfolio-server-flake";
            title = "portfolio-server Backend Flake";
          };
          namespace = "portfolio-server";
        };
      };
    in
    lib.mkFlake {
      channels-config = {
        allowUnfree = false;
      };
      overlays = [ ];
      systems.modules = {
        nixos = [ ];
        darwin = [ ];
      };
      outputs-builder = channels: { formatter = channels.nixpkgs.nixfmt-rfc-style; };
    };
}
