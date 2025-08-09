{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  name = "superhero-api-wrapper";

  buildInputs = with pkgs; [
    go
  ];
}
