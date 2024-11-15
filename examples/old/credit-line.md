```logi macro
macro creditLineRule {
    kind Syntax
    
    types {
        Condition { expr }
    }
    
    syntax {
        Name <name string>
        Description <description string>
        CreditLimit <creditLimit number>
        InterestRate <interestRate number>
        
        Conditions [<condition Condition>]
    }
}
```

```logi
creditLineRule BasicCredit {
    Name "Basic Credit"
    Description "This is a basic credit line."
    CreditLimit 1000
    InterestRate 14
    
    Conditions [
        salary > 300,
        creditScore > 700,
        employmentYears > 2
    ]   
}

creditLineRule PremiumCredit {
    Name "Premium Credit"
    Description "This is a premium credit line."
    CreditLimit 5000
    InterestRate 10
    
    Conditions [
        salary > 500,
        creditScore > 750,
        employmentYears > 3
    ]   
}

creditLineRule GoldCredit {
    Name "Gold Credit"
    Description "This is a gold credit line."
    CreditLimit 10000
    InterestRate 8
    
    Conditions [
        salary > 1000,
        creditScore > 800,
        employmentYears > 5
    ]   
}
```