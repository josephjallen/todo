package todostore

import (
	"context"
	"encoding/json"
	"errors"
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

var lists map[string]*TodoList = make(map[string]*TodoList)

const (
	StatusNotStarted string = "not started"
	StatusStarted    string = "started"
	StatusCompleted  string = "completed"
)

func GetList(ctx context.Context, todoListName string) (*TodoList, error) {
	list, ok := lists[todoListName]
	if !ok {
		logger.InfoLog(ctx, "Init TodoStore for todolist: "+todoListName)
		var err error
		list, err = retrieveListFromFile(ctx, todoListName)
		if err != nil {
			return nil, err
		}
		if list == nil {
			logger.InfoLog(ctx, "Creating todo list: "+todoListName)
			list = &TodoList{Name: todoListName}
		}
		lists[todoListName] = list
		logger.InfoLog(ctx, "Added list to TodoStore: "+todoListName)
	}

	return list, nil
}

func CreateList(ctx context.Context, todoListName string) (*TodoList, error) {
	list, ok := lists[todoListName]
	if ok {
		return list, errors.New("List already exists: " + todoListName)
	}

	list, err := retrieveListFromFile(ctx, todoListName)
	if err != nil {
		return nil, err
	}
	if list != nil {
		return list, errors.New("List already exists: " + todoListName)
	}

	logger.InfoLog(ctx, "Creating todo list: "+todoListName)
	list = &TodoList{Name: todoListName}

	logger.InfoLog(ctx, "Adding list to TodoStore: "+todoListName)
	lists[todoListName] = list

	return list, nil
}

func AddItemToList(ctx context.Context, listName string, itemName string, itemDescription string) error {
	list, err := GetList(ctx, listName)
	if err != nil {
		return err
	}
	var alreadyExists bool = false
	for _, lItem := range list.LItems {
		if lItem.Name == itemName {
			alreadyExists = true
			break
		}
	}

	if alreadyExists {
		return errors.New("Item already exists: " + itemName)
	}

	lItem := TodoListItem{Name: itemName, Description: itemDescription, Status: StatusNotStarted}
	list.LItems = append(list.LItems, lItem)
	logger.InfoLog(ctx, "Added item: "+lItem.Name+" to List: "+list.Name)

	return nil
}

func UpdateListItemDescription(ctx context.Context, listName string, itemName string, itemDescription string) error {
	list, err := GetList(ctx, listName)
	if err != nil {
		return err
	}
	var itemFound bool = false
	var updateItemIndex int
	for index, lItem := range list.LItems {
		if lItem.Name == itemName {
			updateItemIndex = index
			itemFound = true
			break
		}
	}
	if itemFound {
		list.LItems[updateItemIndex].Description = itemDescription
		logger.InfoLog(ctx, "Item Updated (Description): "+itemName+" in List: "+list.Name)
	} else {
		return errors.New("Cannot find Item to update: " + itemName)
	}

	return nil
}

func UpdateListItemStatus(ctx context.Context, listName string, itemName string, itemStatus string) error {
	list, err := GetList(ctx, listName)
	if err != nil {
		return err
	}
	if itemStatus != StatusNotStarted && itemStatus != StatusStarted && itemStatus != StatusCompleted {
		err := errors.New("Invalid status provided: " + itemStatus)
		return err
	}

	var itemFound bool = false
	var updateItemIndex int
	for index, lItem := range list.LItems {
		if lItem.Name == itemName {
			updateItemIndex = index
			itemFound = true
			break
		}
	}
	if itemFound {
		list.LItems[updateItemIndex].Status = itemStatus
		logger.InfoLog(ctx, "Item Updated (Status): "+itemName+" in List: "+list.Name)
	} else {
		return errors.New("Cannot find Item to update: " + itemName)
	}

	return nil
}

func DeleteItemFromList(ctx context.Context, listName string, itemName string) error {
	list, err := GetList(ctx, listName)
	if err != nil {
		return err
	}
	var itemFound bool = false
	var deleteItemIndex int
	for index, lItem := range list.LItems {
		if lItem.Name == itemName {
			deleteItemIndex = index
			itemFound = true
			break
		}
	}
	if itemFound {
		list.LItems = append(list.LItems[:deleteItemIndex], list.LItems[deleteItemIndex+1:]...)
		logger.InfoLog(ctx, "Item Deleted: "+itemName+" from List: "+list.Name)
	} else {
		logger.InfoLog(ctx, "Cannot find Item to delete: "+itemName)
	}

	return nil
}

/*
 */
func retrieveListFromFile(ctx context.Context, listName string) (*TodoList, error) {
	list_b, err := filestorage.LoadFileToByteSlice(ctx, listName+".json")
	if err != nil || list_b == nil {
		return nil, err
	}

	var list TodoList

	if list_b != nil {
		logger.InfoLog(ctx, "Getting todo list: "+listName)
		logger.InfoLog(ctx, string(list_b))
		err := json.Unmarshal(list_b, &list)
		if err != nil {
			return nil, err
		}
	}

	return &list, nil
}

func SaveList(ctx context.Context, listName string) error {
	list, err := GetList(ctx, listName)
	if err != nil {
		return err
	}
	list_bb, err := json.Marshal(list)
	if err != nil {
		return err
	}

	filestorage.SaveByteSliceToFile(list_bb, listName+".json")

	logger.InfoLog(ctx, string(list_bb))

	return nil
}
