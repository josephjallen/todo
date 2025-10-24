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

var list TodoList

func Init(todoListName string, ctx context.Context) error {
	list_, err := getList(todoListName, ctx)
	if err != nil {
		logger.ErrorLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " error:", err)
		return err
	}
	list = list_
	return nil
}

func AddItemToList(itemName string, itemDescription string, ctx context.Context) {
	for _, lItem := range list.LItems {
		if lItem.Name == itemName {
			logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " Item already exists: "+lItem.Name)
		}
	}
	lItem := TodoListItem{Name: itemName, Description: itemDescription}
	list.LItems = append(list.LItems, lItem)
	logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " Added item: "+lItem.Name+" to list: "+list.Name)
}

func UpdateListItem(itemName string, itemDescription string, ctx context.Context) {
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
}

func DeleteItemFromList(itemName string, ctx context.Context) {
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
}

/*
 */
func getList(todoListName string, ctx context.Context) (TodoList, error) {

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
			return TodoList{}, err
		}
	} else {
		list = TodoList{Name: todoListName}
	}

	return list, nil
}

func SaveList(ctx context.Context) error {

	list_bb, err := json.Marshal(list)
	if err != nil {
		logger.ErrorLog.Println(ctx.Value(logger.TraceIdKey{}).(string), " error:", err)
		return err
	}

	filestorage.SaveByteSliceToFile(list_bb, list.Name+".json")

	logger.InfoLog.Println(ctx.Value(logger.TraceIdKey{}).(string), string(list_bb))

	return nil
}
