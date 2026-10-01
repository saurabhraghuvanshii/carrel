An address book has collected the same people several times, from the phone, from mail and from old imports. Each contact has a list of email addresses. Two contacts are the same person if they share an address, and this carries through chains: if contact A shares an address with B, and B shares another with C, then all three are one person, even when A and C have nothing in common.

Contacts are numbered from 0 in the order given. Return one number per contact: the lowest contact number among all contacts that are the same person as it. The app keeps that contact and folds the others into it.

Use union-find over the contacts. Go through the contacts with a hash map from address to the first contact that had it; when an address was seen before, union the two contacts. At the end, the answer for a contact is the lowest number in its set.

## Constraints

- 1 ≤ n ≤ 3000
- 1 ≤ addresses per contact ≤ 5, all different inside one contact
- An address is 1 to 20 characters with no spaces

## Follow-up

The address book syncs all the time: contacts are added, and sometimes an address is removed from a contact, which can split one person back into two. Union-find cannot undo a union. What would you do?
