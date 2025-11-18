package todostore

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"todo/filestorage"
	"todo/logger"
)

type TodoList struct {
	Name   string         `json:"Name"`
	LItems []TodoListItem `json:"lItems"`
}

type TodoListItem struct {
	Name        string
	Description string
	Status      string
}

var List *TodoList
var StatusNotStarted string = "not started"
var StatusStarted string = "started"
var StatusCompleted string = "completed"

func Init(ctx context.Context, todoListName string) error {
	if List == nil {
		filestorage.Init(ctx, todoListName+".json")
		logger.InfoLog(ctx, "Init TodoStore for todolist: "+todoListName)
		var err error
		List, err = getList(ctx, todoListName)
		if err != nil {
			return err
		}
	} else {
		logger.WarningLog(ctx, "TodoStore Single instance already created.")
	}

	return nil
}

func CheckListExists(ctx context.Context, todoListName string) error {
	if List != nil && List.Name == todoListName {
		return nil
	}
	logger.InfoLog(ctx, "Checking if file exists: "+todoListName+".json")
	_, err := os.Stat(todoListName + ".json")
	if os.IsNotExist(err) {
		return errors.New("Todo list does not exist: " + todoListName)
	}

	return nil
}

func AddItemToList(ctx context.Context, itemName string, itemDescription string) error {
	var alreadyExists bool = false
	for _, lItem := range List.LItems {
		if lItem.Name == itemName {
			alreadyExists = true
			break
		}
	}

	if !alreadyExists {
		lItem := TodoListItem{Name: itemName, Description: itemDescription, Status: StatusNotStarted}
		List.LItems = append(List.LItems, lItem)
		logger.InfoLog(ctx, "Added item: "+lItem.Name+" to List: "+List.Name)
	} else {
		return errors.New("Item already exists: " + itemName)
	}

	return nil
}

func UpdateListItemDescription(ctx context.Context, itemName string, itemDescription string) error {
	var updateItemIndex int = -2
	for index, lItem := range List.LItems {
		if lItem.Name == itemName {
			updateItemIndex = index
			break
		}
	}
	if updateItemIndex > -2 {
		List.LItems[updateItemIndex].Description = itemDescription
		logger.InfoLog(ctx, "Item Updated (Description): "+itemName+" in List: "+List.Name)
	} else {
		return errors.New("Cannot find Item to update: " + itemName)
	}

	return nil
}

func UpdateListItemStatus(ctx context.Context, itemName string, itemStatus string) error {

	if itemStatus != StatusNotStarted && itemStatus != StatusStarted && itemStatus != StatusCompleted {
		err := errors.New("Invalid status provided: " + itemStatus)
		return err
	}

	var updateItemIndex int = -2
	for index, lItem := range List.LItems {
		if lItem.Name == itemName {
			updateItemIndex = index
			break
		}
	}
	if updateItemIndex > -2 {
		List.LItems[updateItemIndex].Status = itemStatus
		logger.InfoLog(ctx, "Item Updated (Status): "+itemName+" in List: "+List.Name)
	} else {
		return errors.New("Cannot find Item to update: " + itemName)
	}

	return nil
}

func DeleteItemFromList(ctx context.Context, itemName string) error {
	var deleteItemIndex int = -2
	for index, lItem := range List.LItems {
		if lItem.Name == itemName {
			deleteItemIndex = index
			break
		}
	}
	if deleteItemIndex > -2 {
		List.LItems = append(List.LItems[:deleteItemIndex], List.LItems[deleteItemIndex+1:]...)
		logger.InfoLog(ctx, "Item Deleted: "+itemName+" from List: "+List.Name)
	} else {
		logger.InfoLog(ctx, "Cannot find Item to delete: "+itemName)
	}

	return nil
}

/*
 */
func getList(ctx context.Context, todoListName string) (*TodoList, error) {
	list_b, err := filestorage.LoadFileToByteSlice(ctx)
	if err != nil {
		return &TodoList{}, err
	}

	var list TodoList

	if list_b != nil {
		logger.InfoLog(ctx, "Getting todo list: "+todoListName)
		logger.InfoLog(ctx, string(list_b))
		err := json.Unmarshal(list_b, &list)
		if err != nil {
			return &TodoList{}, err
		}
	} else {
		logger.InfoLog(ctx, "Creating todo list: "+todoListName)
		list = TodoList{Name: todoListName}
	}

	return &list, nil
}

func SaveList(ctx context.Context) error {

	list_bb, err := json.Marshal(List)
	if err != nil {
		return err
	}

	filestorage.SaveByteSliceToFile(list_bb)

	logger.InfoLog(ctx, string(list_bb))

	return nil
}
