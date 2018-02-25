# go-ini

the way i like

ps. this brain cell of perfect soa system is fully completed.

possible upgrades: bindable http api (only?)



sample config:

```
main {
  c1 = w0 w1=w2 w3=w4 
  a1 = true 
  b1 = 0 
  d1 = from program 
}

skaiciai {
  pirmas = 51.3 
  antras = 51 
}

acl {
  10.10.10.10/32 
  10.10.10.11 
}

iplist {
  10.10.10.10 
}

10.10.10.10 {
  users = 10 
  limit = 1M 
}

on connect {
  hook3 
}

notifier {
  token = the_secret 
  10.10.10.20 
  10.10.10.21 
  10.10.10.22 
  10.10.10.23 // just comment
  service_url = http://10.10.10.10:12345/register 
}

```

simple api. supports only: string, int, int64, float32, float64 and bool types.

cfg := ini.New() // now you can have several configs in one app

cfg.Read('config.ini') - parses the file

cfg.Save() - saves config to file

cfg.Write('config.ini') - writes config to file


cfg.Set("section", "key", "value") - adds item to config if value is empty, add key.

cfg.Delete("section", "key") - if key is empty, deletes section


cfg.GetKeysList("section") - returns all keys without values


cfg.Exists("section", "key") - returns true if section has that key

