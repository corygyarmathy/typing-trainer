# NixOS module for deploying typing-trainer to the homelab.
#
# Usage from your existing flake:
#
#   inputs.typing-trainer.url = "github:corygyarmathy/typing-trainer";
#
#   modules = [
#     typing-trainer.nixosModules.default
#     {
#       services.typing-trainer = {
#         enable = true;
#         databaseUrl = config.sops.secrets.typing-trainer.databaseURL.path;
#         jwtSecretFile = config.sops.secrets.typing-trainer.jwt.path;
#       };
#     }
#   ];
#
#   systemd.services.typing-trainer.serviceConfig = {
#   LoadCredential = [
#     "jwt:${cfg.jwtSecretFile}"
#     "dburl:${cfg.databaseUrlFile}"
#   ];
# };
# systemd.services.typing-trainer.environment = {
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
  cfg = config.services.typing-trainer;
in
{
  options.services.typing-trainer = {
    enable = mkEnableOption "typing-trainer backend service";
    # TODO: package, databaseUrl, listenAddress, jwtSecretFile, etc.
  };

  config = mkIf cfg.enable {
    # TODO: systemd.services.typing-trainer = { ... };
  };
}
