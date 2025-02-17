## Interpret aritmetických výrazů

---

**Interpret aritmetických výrazů**. Tyto výrazy obsahují operátory `+, -, \*, /` (s běžnými prioritami a levou asociativitou) a závorky.

Pro zjednodušení úlohy uvažujeme **pouze binární operátory**. V našem jazyce nejsou žádné unární operátory. Navíc můžeme používat **pouze kladná celá čísla** ve výrazech.

### Specifikace vstupu

První řádek vstupu obsahuje číslo `N`. Určuje počet výrazů, které má váš program vyhodnotit. Tyto výrazy jsou na následujících `N` řádcích. Každý řádek obsahuje právě jeden výraz.

### Specifikace výstupu

Pro každý výraz vypište jeden řádek obsahující výsledek – vypočtenou hodnotu výrazu. Pokud je ve vstupu jakákoliv chyba, vypište text `ERROR`.

### Příklad:

**vstup:**

```
2
2 \* (3+5)
15 - 2\*\*7
```

**výstup:**

```
16
ERROR
```

### Spuštění a kompilace:

-   Je třeba mít nainstalovaný **[Golang](https://go.dev/)**
-   Zkompilujte kód pomocí `go build ./main.go`
-   Spusťte zkompilovanou binárku `./main` (nebo `./main.exe`)
-   Otestujte interpret
