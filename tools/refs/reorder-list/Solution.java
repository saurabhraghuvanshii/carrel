// Ported from the owner's linkedlist/LinkedList.java (zigZag). Fixed: the loop tested fast.next twice
// instead of fast != null && fast.next != null, so it crashed; and it cut the
// list at slow.next but reversed from mid.next, losing a node. It now finds
// the middle like the owner's getMid, cuts after it and reverses the rest.
class Solution {
    public void reorder(ListNode head) {
        ListNode slow = head;
        ListNode fast = head.next;
        while (fast != null && fast.next != null) {
            slow = slow.next;
            fast = fast.next.next;
        }
        ListNode curr = slow.next;
        slow.next = null;
        ListNode prev = null;
        ListNode next;
        while (curr != null) {
            next = curr.next;
            curr.next = prev;
            prev = curr;
            curr = next;
        }
        ListNode left = head;
        ListNode right = prev;
        ListNode nextL, nextR;
        while (left != null && right != null) {
            nextL = left.next;
            left.next = right;
            nextR = right.next;
            right.next = nextL;
            left = nextL;
            right = nextR;
        }
    }
}
