package vatsim_webhooks

type Field string

const (
	FieldID               Field = "id"
	FieldNameFirst        Field = "name_first"
	FieldNameLast         Field = "name_last"
	FieldEmail            Field = "email"
	FieldRating           Field = "rating"
	FieldPilotRating      Field = "pilotrating"
	FieldSuspensionDate   Field = "susp_date"
	FieldRegistrationDate Field = "reg_date"
	FieldRegionID         Field = "region_id"
	FieldDivisionID       Field = "division_id"
	FieldSubdivisionID    Field = "subdivision_id"
	FieldLastRatingChange Field = "lastratingchange"
)
