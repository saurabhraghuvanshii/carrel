import java.util.*;

// Ported from the owner's stack/Nextgrt.java, which finds the next greater
// value; here it records the distance to it instead.
class Solution {
    public int[] daysUntilWarmer(int[] arr) {
        Stack<Integer> s = new Stack<>();
        int[] wait = new int[arr.length];
        for (int i = arr.length - 1; i >= 0; i--) {
            while (!s.isEmpty() && arr[s.peek()] <= arr[i]) {
                s.pop();
            }
            wait[i] = s.isEmpty() ? 0 : s.peek() - i;
            s.push(i);
        }
        return wait;
    }
}
