{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  name = "g-browser";

  buildInputs = with pkgs; [
    go
  ];
}
