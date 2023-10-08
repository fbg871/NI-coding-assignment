package main

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const NewestSoftwareVersion = "v3.3.3"

func findAllKomps(ctx context.Context, conn *sqlx.DB) ([]komp, error) {
	var komps []komp
	err := conn.SelectContext(ctx,
		&komps, `
		select id, serial_number, state, software_version, product_code, mac_address, comment
		from Komps
		order by serial_number`,
	)

	if err != nil {
		return nil, err
	}

	var attributes []attribute
	err = conn.SelectContext(
		ctx,
		&attributes,
		`select komp_id, name, value from Attributes order by komp_id, name`,
	)

	if err != nil {
		return nil, err
	}

	for idx, komp := range komps {
		komps[idx].Attributes = make([]attribute, 0)
		for _, attribute := range attributes {
			if komp.ID == attribute.ID {
				komps[idx].Attributes = append(komps[idx].Attributes, attribute)
			}
		}
	}

	return komps, nil
}

func findKompWithSerialNumber(ctx context.Context, serialNumber string, conn *sqlx.DB) (komp, error) {
	var k komp
	err := conn.GetContext(ctx,
		&k, `
		select id, serial_number, state, software_version, product_code, mac_address, comment
		from Komps
		where serial_number = ?`,
		serialNumber,
	)

	if err != nil {
		return komp{}, err
	}

	var attributes = make([]attribute, 0)
	err = conn.SelectContext(
		ctx,
		&attributes,
		`select komp_id, name, value from Attributes where komp_id = ? order by name`,
		k.ID,
	)

	if err != nil {
		return komp{}, err
	}

	k.Attributes = attributes
	return k, nil
}

func updateKompStateAndComment(ctx context.Context, serialNumber, state, comment string, conn *sqlx.DB) (komp, error) {
	k, err := findKompWithSerialNumber(ctx, serialNumber, conn)
	if err != nil {
		return komp{}, err
	}

	if state != "" {
		k.State = state

		if state == "available" && k.SoftwareVersion != NewestSoftwareVersion {
			k.Comment = fmt.Sprintf("Software version is %s. Software upgrade required", k.SoftwareVersion)
		}
	}

	if comment != "" {
		k.Comment = comment
	}

	_, err = conn.NamedExecContext(ctx, `
		update Komps set state=:state, comment=:comment where id=:id
	`, k,
	)

	if err != nil {
		return komp{}, err
	}

	return k, nil
}

func createKompAttribute(ctx context.Context, createdAttribute attribute, conn *sqlx.DB) (attribute, error) {
	fmt.Println("createKompAttribute")

	_, err := conn.NamedExecContext(ctx, 
		`insert into Attributes (komp_id, name, value) values (:komp_id, :name, :value)`, 
		createdAttribute,
	)

	if err != nil {
		return attribute{}, err
	}
	
	return createdAttribute, nil
}

func updateKompAttribute(ctx context.Context, updatedAttribute attribute, conn *sqlx.DB) (attribute, error) {
	fmt.Println("updateKompAttribute")

	_, err := conn.NamedExecContext(ctx, 
		`update Attributes set value=:value where komp_id=:komp_id and name=:name`, 
		updatedAttribute,
	)

	if err != nil {
		return attribute{}, err
	}

	return updatedAttribute, nil
}
