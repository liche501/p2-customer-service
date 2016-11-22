package model

import "time"

type Customer struct{
        Id int64
        Mobile string
        CreatedAt time.Time
        UpdatedAt time.Time
}
