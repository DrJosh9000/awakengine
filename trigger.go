// Copyright 2016 Josh Deprez
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package awakengine

import "github.com/DrJosh9000/vec"

// Trigger is everything to do with reacting to the player or time or ...
// On the PC entering any of the Tiles, Fired, Active, and Depends will be
// checked and then Fire will happen.
// If no Tiles are listed, it will be added to a global list of triggers
// checked on every frame.
type Trigger struct {
	Name    string
	Tiles   []vec.I2
	Active  func(gameFrame int) bool
	Depends []string
	Fire    func(gameFrame int)
	Repeat  bool

	fired bool
}

func (t *Trigger) Reset() { t.fired = false }
