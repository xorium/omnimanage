package converters

import (
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"strconv"
)

func JSONSrcToWeb(src interface{}) (map[string]interface{}, error) {
	s, ok := src.(datatypes.JSON)
	if !ok {
		return nil, fmt.Errorf("Wrong type '%T'", src)
	}

	w := map[string]interface{}{}
	err := json.Unmarshal(s, &w)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func JSONWebToSrc(web interface{}) (datatypes.JSON, error) {
	w, ok := web.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Wrong type '%T'", web)
	}

	j, err := json.Marshal(w)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func IDWebToSrc(web interface{}) (int, error) {
	w, ok := web.(string)
	if !ok {
		return 0, fmt.Errorf("Wrong type '%T'", web)
	}
	id, err := strconv.Atoi(w)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func IDSrcToWeb(src interface{}) (string, error) {
	s, ok := src.(int)
	if !ok {
		return "", fmt.Errorf("Wrong type '%T'", src)
	}
	id := strconv.Itoa(s)
	return id, nil
}
