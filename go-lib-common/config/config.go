package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/constant"
	commonTime "github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
)

// BindFromFile load config from filename then assign to destination
func BindFromFile(dest any, filename string, paths ...string) error {
	v := viper.New()

	v.SetConfigType("json")
	v.SetConfigName(filename)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	for _, path := range paths {
		v.AddConfigPath(path)
	}

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	fmt.Printf("using config file: %s.\n", filename)

	err = v.Unmarshal(dest)
	if err != nil {
		return err
	}

	err = SetEnvFromConsulKV(v)
	if err != nil {
		fmt.Printf("failed to set env from file: %+v\n", err)
		return err
	}

	return nil
}

// BindFromConsul load config from remote consul then assign to destination
func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	err := v.AddRemoteProvider("consul", endPoint, path)
	if err != nil {
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		return err
	}

	fmt.Printf("using config from consul: %s/%s.\n", endPoint, path)

	err = v.Unmarshal(dest)
	if err != nil {
		fmt.Printf("failed to unmarshal config dest: %+v\n", err)
		return err
	}

	err = SetEnvFromConsulKV(v)
	if err != nil {
		fmt.Printf("failed to set env from consul: %+v\n", err)
		return err
	}

	return nil
}

// BindAndWatchFromConsul load and watch config from remote consul then assign to destination
func BindAndWatchFromConsul(dest any, endPoint, path string, interval int) error {
	location, err := time.LoadLocation(commonTime.LoadTimeZoneFromEnv())
	if err != nil {
		fmt.Printf("failed load location: %s\n", location)
	}

	err = BindFromConsul(dest, endPoint, path)
	if err != nil {
		fmt.Printf("failed reloading consul: %+v\n", err)
		return err
	}

	scheduler := gocron.NewScheduler(location)
	_, err = scheduler.Every(interval).Seconds().Do(func() {
		er := BindFromConsul(dest, endPoint, path)
		if er != nil {
			fmt.Printf("failed reloading consul: %+v\n", er)
		}
	})

	if err != nil {
		fmt.Printf("failed run scheduler specify jobFunc: %+v\n", err)
		return err
	}

	scheduler.StartAsync()

	return nil
}

// LoadConsulIntervalFromEnv get interval value for loading config from consul
func LoadConsulIntervalFromEnv() (int, error) {
	fromEnv := os.Getenv(constant.ConsulWatchInterval)
	if len(fromEnv) <= 0 {
		return constant.DefaultLoadConsulInterval, nil
	}

	interval, err := strconv.Atoi(fromEnv)
	if err != nil {
		fmt.Printf("failed convert %s: %+v\n", constant.ConsulWatchInterval, err)
		return 0, err
	}

	return interval, nil
}

func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)
	if err != nil {
		fmt.Printf("failed to unmarshal config env: %+v\n", err)
		return err
	}

	for k, v := range env {
		var (
			valOf = reflect.ValueOf(v)
			val   string
		)

		switch valOf.Kind() {
		case reflect.String:
			val = valOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Uint()))
		case reflect.Float64:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Float32:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		}

		os.Setenv(k, val)
	}

	return nil
}
