/*
   ZAU Single Sign-On
   Copyright (C) 2021  Daniel A. Hawton <daniel@hawton.org>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package middleware

import (
	"fmt"
	"time"

	log "github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
)

func Logger(c *gin.Context) {
	start := time.Now()
	c.Next()
	end := time.Now()
	latency := end.Sub(start)

	log.Category("web").Info(fmt.Sprintf("%s - %s %s - %d \"%s\" %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Writer.Status(), c.Request.UserAgent(), latency))
}
