package todostore

import (
	"context"
	"encoding/json"
	"fmt"
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
		if list == nil {
			fmt.Println("Creating single instance now.")
			var err error
			list, err = getList(ctx, todoListName)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return nil
}

func AddItemToList(ctx context.Context, itemName string, itemDescription string) error {
	for _, lItem := range list.LItems {
		if lItem.Name == itemName {
			logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " Item already exists: "+lItem.Name)
		}
	}
	lItem := TodoListItem{Name: itemName, Description: itemDescription}
	list.LItems = append(list.LItems, lItem)
	logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " Added item: "+lItem.Name+" to list: "+list.Name)

	err := saveList(ctx)
	if err != nil {
		return err
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
		logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " Item Updated: "+itemName+" in list: "+list.Name)
	} else {
		logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " Cannot find Item to update: "+itemName)
	}

	err := saveList(ctx)
	if err != nil {
		return err
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
		logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " Item Deleted: "+itemName+" from list: "+list.Name)
	} else {
		logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " Cannot find Item to delete: "+itemName)
	}

	err := saveList(ctx)
	if err != nil {
		return err
	}

	return nil
}

/*
 */
func getList(ctx context.Context, todoListName string) (*TodoList, error) {
	filename := todoListName
	filename += ".json"
	logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), string(filename))

	list_b, err := filestorage.LoadFileToByteSlice(filename)
	if err != nil {
		logger.ErrorLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " error:", err)
	}

	var list TodoList

	if list_b != nil {
		logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), string(list_b))
		err := json.Unmarshal(list_b, &list)
		if err != nil {
			logger.ErrorLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " error:", err)
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
		logger.ErrorLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " error:", err)
		return err
	}

	filestorage.SaveByteSliceToFile(list_bb, list.Name+".json")

	logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), string(list_bb))

	return nil
}
