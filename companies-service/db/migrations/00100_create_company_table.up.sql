BEGIN;

CREATE TABLE company (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(15) NOT NULL UNIQUE,
    description VARCHAR(3000),
    amount_of_employees INT NOT NULL CHECK (amount_of_employees >= 0),
    registered BOOLEAN NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship'))
);

COMMIT;