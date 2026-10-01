// Ported from the owner's linkedlist/LinkedList.java (deleteNthfromEnd); head is a parameter and the new
// head is returned.
class Solution {
    public ListNode removeFromEnd(ListNode head, int n) {
        int sz = 0;
        ListNode temp = head;
        while (temp != null) {
            temp = temp.next;
            sz++;
        }
        if (n == sz) {
            head = head.next;
            return head;
        }
        int i = 1;
        int iToFind = sz - n;
        ListNode prev = head;
        while (i < iToFind) {
            prev = prev.next;
            i++;
        }
        prev.next = prev.next.next;
        return head;
    }
}
