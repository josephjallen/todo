package todostore

import (
	"encoding/json"
	"todo/filestorage"
	"todo/logger"
)

type TodoList struct {
	Name   string         `json:"name"`
	LItems []TodoListItem `json:"litems"`
}

type TodoListItem struct {
	Name        string
	Description string
}

func AddItemToList(list *TodoList, itemName string, itemDescription string) {
	for _, lItem := range list.LItems {
		if lItem.Name == itemName {
			logger.InfoLog.Println("Item already exists: " + lItem.Name)
			return
		}
	}
	lItem := TodoListItem{Name: itemName, Description: itemDescription}
	list.LItems = append(list.LItems, lItem)
	logger.InfoLog.Println("Added item: " + lItem.Name + " to list: " + list.Name)
}

func UpdateListItem(list *TodoList, itemName string, itemDescription string) {
	var updateItemIndex int = -2
	for index, lItem := range list.LItems {
		if lItem.Name == itemName {
			updateItemIndex = index
			break
		}
	}
	if updateItemIndex > -2 {
		list.LItems[updateItemIndex].Description = itemDescription
		logger.InfoLog.Println("Item Updated: " + itemName + " in list: " + list.Name)
	} else {
		logger.InfoLog.Println("Cannot find Item to update: " + itemName)
	}
}

func DeleteItemFromList(list *TodoList, itemName string) {
	var deleteItemIndex int = -2
	for index, lItem := range list.LItems {
		if lItem.Name == itemName {
			deleteItemIndex = index
			break
		}
	}
	if deleteItemIndex > -2 {
		list.LItems = append(list.LItems[:deleteItemIndex], list.LItems[deleteItemIndex+1:]...)
		logger.InfoLog.Println("Item Deleted: " + itemName + " from list: " + list.Name)
	} else {
		logger.InfoLog.Println("Cannot find Item to delete: " + itemName)
	}
}

/*
 */
func GetList(todoListName string) (*TodoList, error) {

	filename := todoListName
	filename += ".json"
	logger.InfoLog.Println(string(filename))

	list_b, err := filestorage.LoadFileToByteSlice(filename)
	if err != nil {
		logger.ErrorLog.Println("error:", err)
	}

	var list TodoList

	if list_b != nil {
		logger.InfoLog.Println(string(list_b))
		err := json.Unmarshal(list_b, &list)
		if err != nil {
			logger.ErrorLog.Println("error:", err)
			return nil, err
		}
	} else {
		list = TodoList{Name: todoListName}
	}

	return &list, nil
}

func SaveList(list *TodoList) error {

	list_bb, err := json.Marshal(list)
	if err != nil {
		logger.ErrorLog.Println("error:", err)
		return err
	}

	filestorage.SaveByteSliceToFile(list_bb, list.Name+".json")

	logger.InfoLog.Println(string(list_bb))

	return nil
}
