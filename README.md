<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# terraform-provider-satellite

A Terraform provider for Red Hat Satellite.

This provider can be used to create and manage organizations, subscription manifests, roles, activation keys,
and host collections among other things. Generally speaking, things that are useful to implement Satellite in
a large multitennant environment with or without multiple organizations.

The provider does not have working tests so it should probably be considered beta.

## Building/Installing

The provider is now available on the [Terraform Registry](https://registry.terraform.io/providers/umich-vci/satellite/latest)
so you probably don't need to build the provider unless you want to contribute.

That said, running `GO111MODULE=on go get -u github.com/umich-vci/terraform-provider-satellite` should download
the code and result in a binary at `$GOPATH/bin/terraform-provider-satellite`. You can then move the
binary to `~/.terraform.d/plugins` to use it with Terraform.

This has been tested with Terraform 1.3.x and Satellite 6.11.x. as of version 0.7.0. Previous versions of Satellite may not work
with this version as there have been API changes that are difficult to track.

## License

This project is licensed under the Mozilla Public License Version 2.0.
