ListMerchantOmzet := `
    select 
		t.merchant_id,
		m.merchant_name,
		IF(sum(t.bill_total)=null, 0, sum(t.bill_total)) as omzet,
		DATE_FORMAT(t.created_at , '%Y-%m-%d') as created_at 
	from 
		transactions t
	left join 
		merchants m on m.id = t.merchant_id 
	left join 
		users u on u.id = m.user_id
	where 
		DATE_FORMAT(t.created_at , '%Y-%m-%d') = ?
	and
		u.user_name = ?
	group by t.merchant_id, created_at
	order by t.merchant_id
`

GetListMerchant := `
    select m.id as merchant_id, m.merchant_name from merchants m
	left join users u on u.id = m.user_id
	where u.user_name = ?
`

ListOutletOmzet := `
    select 
			t.merchant_id,
			o.id as outlet_id ,
			m.merchant_name,
			o.outlet_name,
			IF(sum(t.bill_total)=null, 0, sum(t.bill_total)) as omzet,
			DATE_FORMAT(t.created_at , '%Y-%m-%d') as created_at 
		from 
			transactions t
		left join 
			merchants m on m.id = t.merchant_id
		left join
			outlets o on o.id = t.outlet_id 
		left join
			users u on u.id = m.user_id
		where 
			DATE_FORMAT(t.created_at , '%Y-%m-%d') = ?
		and
			u.user_name = ?
		group by t.merchant_id, o.id, created_at
		order by t.merchant_id
`

GetListOutlet := `
    select o.id as outlet_id, o.merchant_id, m.merchant_name ,outlet_name 
	from outlets o
	left join merchants m on m.id = o.merchant_id
	left join users u on u.id = m.user_id
	where u.user_name = ?
	order by o.merchant_id
`

GetDataByUsername := `
    select u.id as user_id,u.name,u.password,u.user_name, u.created_at, u.updated_at 
	from users u
	where u.user_name = ?
	limit 1
`