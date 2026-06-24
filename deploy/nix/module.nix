# NixOS module for deploying typist to the homelab.
#
# Usage from your existing flake:
#
#   inputs.typist.url = "github:corygyarmathy/typist";
#
#   modules = [
#     typist.nixosModules.default
#     {
#       services.typist = {
#         enable = true;
#         databaseUrl = config.sops.secrets.typist.databaseURL.path;
#         jwtSecretFile = config.sops.secrets.typist.jwt.path;
#       };
#     }
#   ];
#
#   systemd.services.typist.serviceConfig = {
#   LoadCredential = [
#     "jwt:${cfg.jwtSecretFile}"
#     "dburl:${cfg.databaseUrlFile}"
#   ];
# };
# systemd.services.typist.environment = {
#   JWT_SECRET_FILE    = "%d/jwt";    # %d = $CREDENTIALS_DIRECTORY
#   DATABASE_URL_FILE  = "%d/dburl";
# };
#
# TODO(phase-7): implement options, systemd unit, postgresql provisioning.

{
  config,
  lib,
  pkgs,
  ...
}:

with lib;

let
  cfg = config.services.typist;
in
{
  options.services.typist = {
    enable = mkEnableOption "typist backend service";
    # TODO: package, databaseUrl, listenAddress, jwtSecretFile, etc.
  };

  config = mkIf cfg.enable {
    # TODO: systemd.services.typist = { ... };
  };
}
