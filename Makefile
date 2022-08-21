dot:
	go run golex.go > ex.dot; dot -Kdot -Tpng ex.dot -oex.png
