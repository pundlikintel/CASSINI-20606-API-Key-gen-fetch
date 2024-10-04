Running tool

Step 1. 
    Build Tool
```    
    go build -o partnerkeygen .
```
Step 2.
    Run tool
```
    ./partnerkeygen -issuer="<value>" -sec_kid="<value>" -aws_key_id="<value>" -aws_secret_key="<value>" -aws_token="<value>" -region="<value>" -tenant_id="<value>" -user_id="<value>" -environment="<value>" -product_id="<value>" -service_id="<value>"
```
  
| Flag           | Meaning                                                    |
|----------------|------------------------------------------------------------|
| issuer         | domain value like "api-dev01-user7.project-amber-smas.com" |
| sec_kid        | secret id                                                  |
| aws_key_id     | aws key id                                                 |
| aws_secret_key | aws secret key                                             |
| aws_token      | aws token                                                  |                        
| region         | aws region                                                 |                       
| tenant_id      | parter tenant id                                           |                     
| user_id        | partener user id                                           |                     
| environment    | full environment name like "dev01-user7"                   | 
| product_id     | product id                                                 | 
| service_id     | Service id                                                 |

