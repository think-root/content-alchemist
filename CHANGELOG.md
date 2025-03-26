# [1.6.0](https://github.com/think-root/content-alchemist/compare/v1.5.0...v1.6.0) (2025-03-26)


### Bug Fixes

* change Posted field to pointer in getRepositoryRequestBody ([f60fe2c](https://github.com/think-root/content-alchemist/commit/f60fe2c156654a6029d9f36e912a3a11820602c3))


### Features

* implement pagination in GetRepository function ([543ae86](https://github.com/think-root/content-alchemist/commit/543ae8626baf023fba51060ee0f953cc1e52c86d))

# [1.5.0](https://github.com/think-root/content-alchemist/compare/v1.4.0...v1.5.0) (2025-03-24)


### Features

* add CORS headers to API responses ([428f6b4](https://github.com/think-root/content-alchemist/commit/428f6b4840fe29f23907fb291f2ba8b07c29e349))
* add CORS middleware ([589a2b1](https://github.com/think-root/content-alchemist/commit/589a2b1066fca05786e8f37aadc492b9dcba8358))
* add CORS middleware to server ([c74aabf](https://github.com/think-root/content-alchemist/commit/c74aabf300a96934b2d186aecf899b3794c67279))

# [1.4.0](https://github.com/think-root/content-alchemist/compare/v1.3.0...v1.4.0) (2025-03-21)


### Features

* add DateAdded field to repository entry in database ([7480505](https://github.com/think-root/content-alchemist/commit/7480505fca26cae409ad6d01f82e1686bc9c0eb0))
* add DatePosted field to GithubRepositories model for enhanced date tracking ([66f9edd](https://github.com/think-root/content-alchemist/commit/66f9eddeeba51441b5d8f735e1cce2094bb74f1f))
* enhance GetRepository function to support sorting by date and include DateAdded and DatePosted in response ([e1768d8](https://github.com/think-root/content-alchemist/commit/e1768d859b985ef79a63697f8f5bd8d181ec4edd))
* enhance GetRepository function to support sorting by date posted, date added, or id with order options ([81c3156](https://github.com/think-root/content-alchemist/commit/81c3156c06719382dec2fda93d09598c57a33126))
* update UpdatePostedStatusByURL to set DatePosted when repository is marked as posted ([99f6c68](https://github.com/think-root/content-alchemist/commit/99f6c68a97e918653916b22231286f45fac62b3d))

# [1.3.0](https://github.com/think-root/content-alchemist/compare/v1.2.0...v1.3.0) (2025-02-22)

### Features

- add APP_VERSION argument to Docker Compose configuration ([8cf4056](https://github.com/think-root/content-alchemist/commit/8cf40569af11c82d2fb1a8b13b5f9cd48a076d27))
- add APP_VERSION argument to Dockerfile for build configuration ([b7c3fc8](https://github.com/think-root/content-alchemist/commit/b7c3fc8d5fbf58972b9ae981870cf68003e39b3b))
- update deployment workflow to determine APP_VERSION from Git tags ([b568d84](https://github.com/think-root/content-alchemist/commit/b568d842a88429b24de253cbb91715b445e66de4))

# [1.2.0](https://github.com/think-root/content-alchemist/compare/v1.1.0...v1.2.0) (2025-02-17)

### Features

- add EnvAsInt function to retrieve environment variables as integers ([4dd8c32](https://github.com/think-root/content-alchemist/commit/4dd8c32db7da2206f250f79b31fb7c1bd6537d77))
- add rate limiting middleware to server ([bf57eeb](https://github.com/think-root/content-alchemist/commit/bf57eeba077ab31e2530b2145526afaf383676e6))
- add RATE_LIMIT configuration variable to config ([160c9d0](https://github.com/think-root/content-alchemist/commit/160c9d0a0000049bd8bb785f172553edc59dc7a0))
- implement rate limiting middleware for HTTP requests ([0da1e74](https://github.com/think-root/content-alchemist/commit/0da1e74e36b3e72c161bb64b5e4001e43323f68c))

# [1.1.0](https://github.com/think-root/content-alchemist/compare/v1.0.1...v1.1.0) (2025-02-17)

### Bug Fixes

- **workflow:** change npm ci to npm install for dependency installation ([f2a1f1a](https://github.com/think-root/content-alchemist/commit/f2a1f1a053d6d3e40066b0e3ec64c24ca64fcea6))

### Features

- **issue-templates:** create new issue templates for bug reports, feature requests, questions, and documentation ([31ae382](https://github.com/think-root/content-alchemist/commit/31ae382ff3c34adab57be05a737499d491e52fb0))
- **package:** add semantic release dependencies for automated versioning and changelog management ([6e2e0ac](https://github.com/think-root/content-alchemist/commit/6e2e0ac2b8a8e0297dcdf7b4126b55b43fff1c9d))
- **pr-template:** add pull request template for better contribution guidelines ([3306934](https://github.com/think-root/content-alchemist/commit/3306934976935d5664cf50f008a5f8a1bafbdce4))
- **release-config:** add repository URL to semantic release configuration ([0c5bb19](https://github.com/think-root/content-alchemist/commit/0c5bb1980d22e30f3a2bf58164a9fe26c2992285))
- **release-config:** add semantic release configuration for automated versioning and changelog generation ([dbcfe82](https://github.com/think-root/content-alchemist/commit/dbcfe826513705db05c7b8e25850be6887cb728f))
- **templates:** add issue templates for bug reports, feature requests, questions, and documentation ([bdacc14](https://github.com/think-root/content-alchemist/commit/bdacc14f1ce389ef930e4f6336f76f1dd16e4442))
- **workflow:** enhance GitHub Actions configuration for improved deployment and Git setup ([5e1f63d](https://github.com/think-root/content-alchemist/commit/5e1f63db9f618f57c9485b3ca1f78506b55756b2))
- **workflow:** update GitHub Actions for semantic release and improve deployment process ([3449118](https://github.com/think-root/content-alchemist/commit/344911896f7626d6b6bba7519fc569bc979955f0))
