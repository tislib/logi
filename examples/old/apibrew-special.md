ApiBrew Special examples
========================

```logi
resource Indicator {
    id: int [id, required]
    name: string [required]
    description: string
    type: string [required]
    value: string [required]
    created_at: datetime [required, default: now()]
    updated_at: datetime [required, default: now()]
}

code IndicatorCheck {
    beforeCreate(indicator Indicator) {
        if (indicator.type == "string") {
            indicator.value = "string"
        }
    }
}

code IndicatorCheck2 {
    useResource(Indicator2) as Indicator2

    beforeCreate(indicator Indicator) {
       Indicator2.create({name: indicator.name, description: indicator.description, type: indicator.type, value: indicator.value})
    }
}
```