package todostore

import (
	"context"
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

var list *TodoList

func Init(ctx context.Context, todoListName string) error {
	if list == nil {
		filestorage.Init(ctx, todoListName+".json")
		logger.InfoLog(ctx, "TodoStore Creating single instance now.")
		var err error
		list, err = getList(ctx, todoListName)
		if err != nil {
			return err
		}
	} else {
		logger.WarningLog(ctx, "TodoStore Single instance already created.")
	}

	return nil
}

func AddItemToList(ctx context.Context, itemName string, itemDescription string) error {
	var alreadyExists bool = false
	for _, lItem := range list.LItems {
		if lItem.Name == itemName {
			alreadyExists = true
			break
		}
	}

	if !alreadyExists {
		lItem := TodoListItem{Name: itemName, Description: itemDescription}
		list.LItems = append(list.LItems, lItem)
		logger.InfoLog(ctx, "Added item: "+lItem.Name+" to list: "+list.Name)

		err := saveList(ctx)
		if err != nil {
			return err
		}
	} else {
		logger.InfoLog(ctx, "Item already exists: "+itemName)
	}

	return nil
}

func UpdateListItem(ctx context.Context, itemName string, itemDescription string) error {
	var updateItemIndex int = -2
	for index, lItem := range list.LItems {
		if lItem.Name == itemName {
			updateItemIndex = index
			break
		}
	}
	if updateItemIndex > -2 {
		list.LItems[updateItemIndex].Description = itemDescription
		logger.InfoLog(ctx, "Item Updated: "+itemName+" in list: "+list.Name)
		err := saveList(ctx)
		if err != nil {
			return err
		}
	} else {
		logger.InfoLog(ctx, "Cannot find Item to update: "+itemName)
	}

	return nil
}

func DeleteItemFromList(ctx context.Context, itemName string) error {
	var deleteItemIndex int = -2
	for index, lItem := range list.LItems {
		if lItem.Name == itemName {
			deleteItemIndex = index
			break
		}
	}
	if deleteItemIndex > -2 {
		list.LItems = append(list.LItems[:deleteItemIndex], list.LItems[deleteItemIndex+1:]...)
		logger.InfoLog(ctx, "Item Deleted: "+itemName+" from list: "+list.Name)
		err := saveList(ctx)
		if err != nil {
			return err
		}
	} else {
		logger.InfoLog(ctx, "Cannot find Item to delete: "+itemName)
	}

	return nil
}

/*
 */
func getList(ctx context.Context, todoListName string) (*TodoList, error) {
	logger.InfoLog(ctx, string(todoListName))

	list_b, err := filestorage.LoadFileToByteSlice()
	if err != nil {
		return &TodoList{}, err
	}

	var list TodoList

	if list_b != nil {
		logger.InfoLog(ctx, string(list_b))
		err := json.Unmarshal(list_b, &list)
		if err != nil {
			return &TodoList{}, err
		}
	} else {
		list = TodoList{Name: todoListName}
	}

	return &list, nil
}

func saveList(ctx context.Context) error {

	list_bb, err := json.Marshal(list)
	if err != nil {
		return err
	}

	filestorage.SaveByteSliceToFile(list_bb)

	logger.InfoLog(ctx, string(list_bb))

	return nil
}
