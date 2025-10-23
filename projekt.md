# Aplikacja do tworzenia i szukania event'ów

Aplikacja pozwalająca użytkownikom tworzyć eventy, znajdywać je oraz się na nie zapisywać.

## 1. Technologia użyta
Do backendu: *Go*

Frontend: ~~??, pewnie framework do JS'a w stylu react~~ szkic w [bibliotece wbudowanej w Go - http/template](https://pkg.go.dev/html/template). Mam nadziję, że na moje potrzeby będzie wystarczało, jeśli nie to migracja do Reacta.

DB framework: [database/sql](https://pkg.go.dev/database/sql)

## 2. MVP features:
 - Logowanie do systemu
 - 2 poziomy dostępu do systemu: *Host*, *Attendee*
 - Jeśli ktoś ma konto hosta możliwa zmiana widoku `Attendee <=> Host`
 - Dla Hostów:
   - CRUD na eventach (swoich)
   - przeglądanie uczestników w swoich eventach (ilość i imiona, nazwiska)

 - Dla Attendee:
   - przeglądanie eventów
   - odczytywanie informacji o evencie 
   - zapisywanie się na eventy
   - wyszukiwanie eventów po: nazwie, tagach

## 3. Featury kolejnych wersji:
 - Podstawowe manipulacje kontem: zmiana hasła, usunięcie konta
 - Rejestrowanie danych BIO 
 - Dla Hostów:
   - Przeglądanie demografii użytkowników

 - Dla Attendee:
   - wyszukiwanie eventów po lokalizacji
   - podanie zainteresowań => możliwość sugerowania eventów w kręgu zainteresowań
   - Komentarze, oceny eventów

## 4. Modele
 - User{id: uuid, email: str, password hash: str, name: str, role: fk(role), created_at: date}
 - Role{id: uuid, name}
 - Event{id: uuid, name: str, desc: str, date: date, long: decimal, lat: decimal, fee: decimal(nullable), organiser: fk(user)}
 - Tag{id: uuid, name}