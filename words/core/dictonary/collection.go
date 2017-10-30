package dictonary

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Первоначальный размер словаря
const dictonarySize = 46000

// 0 байт
type item struct{}

// Сollection словарь. Не потокобезопасен
type Сollection struct {
	data map[string]item
}

// FromReader загружает данные словаря из потока, каждый новый элемент в новой строке
func (c *Сollection) FromReader(r io.Reader) error {
	c.Reset()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		c.data[scanner.Text()] = item{}
	}

	return scanner.Err()
}

// FromFile загружает данные словаря из файла, каждый новый элемент в новой строке
func (c *Сollection) FromFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("Open dictonary file: %v", err)
	}
	defer f.Close()

	return c.FromReader(f)
}

// Reset обнуляет словарь
func (c *Сollection) Reset() {
	c.data = make(map[string]item, dictonarySize)
}

// HasKey проверяет есть ли ключ в словаре
func (c *Сollection) HasKey(key string) bool {
	_, ok := c.data[key]
	return ok
}

// SetKey устанавливает ключ в словарь
func (c *Сollection) SetKey(key string) {
	c.data[key] = item{}
}

// DeleteKey удаляет ключ в словаря
func (c *Сollection) DeleteKey(key string) {
	delete(c.data, key)
}

// New создаёт новую инициализированую коллекцию
func New() Сollection {
	c := Сollection{}
	c.Reset()
	return c
}
