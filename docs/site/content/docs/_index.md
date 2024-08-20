---
title: Documentation
toc: false
---

`pkgr` is a command-line interface for managing R packages.  This
documentation describes that [interface][commands] and the
[configuration file][config] that controls its operation.

{{< cards >}}
  {{< card link="config" title="Configuration" icon="document-text"
    subtitle="Define the package library to create" >}}
{{< /cards >}}

{{< cards >}}
  {{< card link="commands/plan" title="Key subcommand: plan" icon="terminal"
    subtitle="Preview what would be installed or updated" >}}
  {{< card link="commands/install" title="Key subcommand: install" icon="terminal"
    subtitle="Create or update the library" >}}
{{< /cards >}}

[commands]: {{< relref "/docs/commands" >}}
[config]: {{< relref "/docs/config" >}}
[install]: {{< relref "/docs/commands/install" >}}
[plan]: {{< relref "/docs/commands/plan" >}}
