package packets

// func TestConvertBytes2UpdateMessage(t *testing.T) {
// 	as := 64513
// 	ip := "10.0.100.3"
// 	localAS := 64514
// 	localIP := "10.200.100.3"

// 	updateMessagePathAttribute := []PathAttribute{
// 		pa.Origin(origin.IGP),
// 		pa.AsPath([]AsPath{as, localAS}),
// 		pa.NextHop(localIP),
// 	}

// 	updateMessage := NewUpdateMessage(
// 		updateMessagePathAttribute,
// 		["10.100.220.0/24"],
// 		[],
// 	)

// 	updateMessageBytes := updateMessage.clone().into()

// 	updateMessage2 := updateMessageBytes.try_into()

// 	assert.Equal(updateMessage, updateMessage2)
// }
