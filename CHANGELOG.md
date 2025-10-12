# [2.0.0](https://github.com/think-root/content-alchemist/compare/v1.17.0...v2.0.0) (2025-10-12)


### Bug Fixes

* **router:** return raw repository text if language parameter is missing ([8e4ed5c](https://github.com/think-root/content-alchemist/commit/8e4ed5c40ec4786b050f1d29f415579c53467a97))


### Features

* **router:** add optional text_language field for multilingual updates ([958468e](https://github.com/think-root/content-alchemist/commit/958468e9f420fe81f50d612e00ff16f49aac0f6a))


### BREAKING CHANGES

* **router:** callers must provide the language in the request body
using `text_language`; the previous `lang` query parameter is no longer
supported.

# [1.17.0](https://github.com/think-root/content-alchemist/compare/v1.16.0...v1.17.0) (2025-09-06)


### Features

* **llm:** add support for Chutes API ([753db1a](https://github.com/think-root/content-alchemist/commit/753db1a1e4fba6ee85c5f163805ea7fb4c8886d1))
* **providers:** add support for Chutes.ai ([c688f1e](https://github.com/think-root/content-alchemist/commit/c688f1eadbed45ca12c76ae5080ea39021fd2d76))

# [1.16.0](https://github.com/think-root/content-alchemist/compare/v1.15.0...v1.16.0) (2025-08-04)


### Features

* **server:** add multilingual text cleaning and validation improvements ([f59eaa4](https://github.com/think-root/content-alchemist/commit/f59eaa4e421a653221ce41d493425006ac70fdfb))
* **server:** clean multilingual text before database insertion ([1e94a55](https://github.com/think-root/content-alchemist/commit/1e94a5592248e566aa727007ceb4f1f46abca623))
* **server:** clean multilingual text before database insertion ([7ae25d3](https://github.com/think-root/content-alchemist/commit/7ae25d3c473fa35298f69c1d0bef324c559dae78))

# [1.15.0](https://github.com/think-root/content-alchemist/compare/v1.14.0...v1.15.0) (2025-07-14)


### Features

* **server:** add text cleaning for multilingual content ([e79671e](https://github.com/think-root/content-alchemist/commit/e79671e333b4357291190aaa4a7fafbeec9be0a1))

# [1.14.0](https://github.com/think-root/content-alchemist/compare/v1.13.1...v1.14.0) (2025-07-05)


### Features

* **db:** add function to get repository by ID or URL ([9560b70](https://github.com/think-root/content-alchemist/commit/9560b701e403385668f45c80e3a829e7203aa0a2))
* **server:** add language validator functionality ([a6bebd0](https://github.com/think-root/content-alchemist/commit/a6bebd0e3b9f3c2b4c4a5a24ff31c6f4bdfde469))
* **server:** add multilingual helper functions ([894462f](https://github.com/think-root/content-alchemist/commit/894462f7970f03a4b4995c274b1184e706361da0))
* **server:** add multilingual support to auto generate endpoint ([bd6acc1](https://github.com/think-root/content-alchemist/commit/bd6acc182115ed7113fe22512628c88cf13ab26d))
* **server:** add multilingual support to manual generate endpoint ([ca3ac58](https://github.com/think-root/content-alchemist/commit/ca3ac58efd1e0de98967748fa088aa01985ba15f))
* **server:** add multilingual support to update repository endpoint ([ed4ace5](https://github.com/think-root/content-alchemist/commit/ed4ace5b339d5c143d379ef348b7a5a7c2006b97))
* **server:** add multilingual text processing to get repository endpoint ([e19fd0a](https://github.com/think-root/content-alchemist/commit/e19fd0af4d0ae024eec6ab7bed95ce263dffb3d6))

## [1.13.1](https://github.com/think-root/content-alchemist/compare/v1.13.0...v1.13.1) (2025-07-04)


### Bug Fixes

* **middlewares:** allow PATCH method in CORS headers ([a93c689](https://github.com/think-root/content-alchemist/commit/a93c689eeb4b524bc486a64b2a8b15f919ade888))

# [1.13.0](https://github.com/think-root/content-alchemist/compare/v1.12.0...v1.13.0) (2025-07-04)


### Features

* **api:** add delete repository endpoint ([9cb35c8](https://github.com/think-root/content-alchemist/commit/9cb35c876cc498fa802ba31c65cbdde7770fb308))
* **api:** add delete repository route to server ([1f5880b](https://github.com/think-root/content-alchemist/commit/1f5880bbd18694e443cbde7108df0a6558ac553d))
* **db:** add delete repository by id or url function ([3215a12](https://github.com/think-root/content-alchemist/commit/3215a12aefa24d357442269176a580f7f4ac5d58))

# [1.12.0](https://github.com/think-root/content-alchemist/compare/v1.11.0...v1.12.0) (2025-06-28)


### Features

* **parser:** add repository filtering logic and database integration ([9af8e7a](https://github.com/think-root/content-alchemist/commit/9af8e7a6d36f7f683dec38ff66e24a67b2dbbccc))

# [1.11.0](https://github.com/think-root/content-alchemist/compare/v1.10.0...v1.11.0) (2025-06-22)


### Features

* **api:** add update repository text endpoint ([98d38de](https://github.com/think-root/content-alchemist/commit/98d38deb099fbd563573de5c4c18d854dcc868fb))
* **db:** add update repository text functionality ([b56555e](https://github.com/think-root/content-alchemist/commit/b56555e98209fde25eb5a0550f82047634eb98cf))

# [1.10.0](https://github.com/think-root/content-alchemist/compare/v1.9.1...v1.10.0) (2025-06-17)


### Features

* **db:** Add unique constraint to repository URLs ([30244ee](https://github.com/think-root/content-alchemist/commit/30244ee6c48cfc0c18f47cf170639037699c2909))

## [1.9.1](https://github.com/think-root/content-alchemist/compare/v1.9.0...v1.9.1) (2025-06-15)


### Bug Fixes

* **db:** include database name in default DSN ([9fdbb5d](https://github.com/think-root/content-alchemist/commit/9fdbb5d967a8b310bafcf9855a1b6b47160d3cd1))

# [1.9.0](https://github.com/think-root/content-alchemist/compare/v1.8.1...v1.9.0) (2025-06-15)


### Bug Fixes

* **server:** handle server startup errors ([763a9dd](https://github.com/think-root/content-alchemist/commit/763a9ddf17837c86116c06e690c84288497ac04a))


### Features

* **ci:** introduce separate Docker Compose files for app and db ([ead3bdd](https://github.com/think-root/content-alchemist/commit/ead3bdd05ce8afc0d6578878922ed46c177cb7fc))
* **db:** migrate database from MySQL to PostgreSQL ([0a40225](https://github.com/think-root/content-alchemist/commit/0a40225ea8db4b1ecf26d0121d5265eb887f7a5b))
* **db:** switch database operations from GORM to raw SQL ([4610c70](https://github.com/think-root/content-alchemist/commit/4610c70c216d05f53bb212dd93ddeb92e4153de5))

## [1.8.1](https://github.com/think-root/content-alchemist/compare/v1.8.0...v1.8.1) (2025-05-27)


### Bug Fixes

* **llm:** correct HTTP_REFERER header in OpenRouter API requests ([8073cf9](https://github.com/think-root/content-alchemist/commit/8073cf9db6a636581b1a19b8726bfac68be54fe1))

# [1.8.0](https://github.com/think-root/content-alchemist/compare/v1.7.0...v1.8.0) (2025-05-07)


### Features

* **api:** add error messages to auto and manual generation responses ([fc95a39](https://github.com/think-root/content-alchemist/commit/fc95a39e20ee93e0bf4dab775c392f1b32a3ba9d))
* **llm:** enhance AutoGenerate request with LLMProvider and configuration options ([50064fd](https://github.com/think-root/content-alchemist/commit/50064fd6008d7da23ef901fd53022ce04f409d9d))
* **llm:** enhance manual generation with LLMProvider and configuration options ([3db14ec](https://github.com/think-root/content-alchemist/commit/3db14ec74a41f49363d91b9b791b2a3072f6309a))
* **llm:** implement LLMProvider interface and message processing functions ([f6ae95e](https://github.com/think-root/content-alchemist/commit/f6ae95ea6affac6fcaeb6a0eaf111f95d282f55d))
* **llm:** implement MistralAgent function for API interaction ([54ff465](https://github.com/think-root/content-alchemist/commit/54ff4657113c9e97d020a18fc872dc68ce66027b))
* **llm:** implement MistralAPI function for API interaction ([a43db3b](https://github.com/think-root/content-alchemist/commit/a43db3b4868a21c3fdb2da7aca747a3086448f92))
* **llm:** implement OpenAI function for API interaction ([5409c67](https://github.com/think-root/content-alchemist/commit/5409c679d275355a7d547bdf7af8e41a9f5d45a3))
* **llm:** implement OpenRouter function for API interaction ([94c59c0](https://github.com/think-root/content-alchemist/commit/94c59c058357e0ed59edc7d5304d3df4be21b650))

# [1.7.0](https://github.com/think-root/content-alchemist/compare/v1.6.0...v1.7.0) (2025-04-03)


### Features

* add indexes and optimize database connection ([3dd2cf2](https://github.com/think-root/content-alchemist/commit/3dd2cf2fb949d4bb50c0eeecdb3fb8d55fb42560))
* improve database query construction in GetRepository ([c911719](https://github.com/think-root/content-alchemist/commit/c911719dbc196c98af6b8089432811d97032b9fa))

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
