# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Introduced a Kubernetes Operator for managing PlantD resources in a Kubernetes cluster. This Operator simplifies the deployment and management of PlantD resources by leveraging Kubernetes native constructs.

### Changed
- PlantD resources can now be defined and managed using Kubernetes Custom Resource Definitions (CRDs). This aligns PlantD with Kubernetes best practices for resource management.
- PlantD configuration can be specified as YAML or JSON manifests, allowing for greater flexibility and compatibility with Kubernetes declarative object syntax.

### Fixed
- Fixed an issue where PlantD resource definitions were not being properly reconciled when changes were made in a Kubernetes environment. The Operator now ensures that PlantD resources are synchronized with the desired state.

### Deprecated
- The traditional method of configuring PlantD resources outside of Kubernetes using standalone configuration files is now deprecated in favor of using the Kubernetes Operator. Future releases may phase out support for standalone configurations.

### Removed
- Removed legacy methods of managing PlantD resources outside of Kubernetes, including traditional configuration files and manual setup scripts. Users are encouraged to migrate to the Kubernetes Operator for resource management.

### Security
- Enhanced security by leveraging Kubernetes RBAC (Role-Based Access Control) for fine-grained access control to PlantD resources within the cluster.

### Documentation
- Updated documentation to include detailed instructions on deploying and configuring PlantD using the Kubernetes Operator.