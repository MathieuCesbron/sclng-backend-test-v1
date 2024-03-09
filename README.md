# Scalingo technical test

## Quickstart
1. Create a [github token](https://docs.github.com/en/enterprise-server@3.9/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token) and export it
```bash
export GITHUB_TOKEN=<your github token>
```
2. Deploy the app
```bash
docker compose up
```
3. The app is running on port `8080`. http://localhost:8080/health should return `OK`.

## API documentation

<details>
 <summary><code>GET</code> <code><b>/health</b></code> <code>app healthcheck</code></summary>

 ##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `text/plain;charset=UTF-8`        | `OK`                                                                |

##### Example curl

> ```bash
>  curl "http://localhost:8080/health"
> ```

</details>

<details>
 <summary><code>GET</code> <code><b>/repos</b></code> <code>get the 100 latest updated public github repos</code></summary>

 ##### Parameters

> | name      |  type     | data type               | description                                                           |
> |-----------|-----------|-------------------------|-----------------------------------------------------------------------|
> | languages      |  optional | string   | Commas separated list of languages. e.g. "python,c". Only repos with at least one language will be returned. Not case-sensitive. You can find the list of supported languages [here](https://github.com/github-linguist/linguist/blob/master/lib/linguist/languages.yml).  |
> | licenses       |  optional | string   | Commas separated list of licenses. e.g. "mit,gpl-3.0". Only repos with one of the license will be returned. Not case-sensitive. You can find the list of supported licenses [here](https://github.com/spdx/license-list-data/blob/main/licenses.md).   |

 ##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json`        | [{"full_name":"maestraccio/team","owner":"maestraccio","repository":"team","license":"gpl-3.0","languages":{"Python":117108}}, ...}]                                                                |

##### Example cURL

> ```bash
>  curl -H "Content-Type: application/json" "http://localhost:8080/repos"
>  curl -H "Content-Type: application/json" "http://localhost:8080/repos?languages=go,python&licenses=mit,gpu-2.0"
> ```

</details>

## How is it scalable ?
- The image size is < 11MB.
- Could scale horizontally on kubernetes.
- No goroutine leaks, canceling the request cancel all goroutines launched.
