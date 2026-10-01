import java.util.*;

// Ported from the owner's Heap/CoonecetNRopes.java: the body of main, now
// returning the cost instead of printing it.
class Solution {
    public int joinCost(int[] ropes) {
        PriorityQueue<Integer> pq = new PriorityQueue<>();
        for (int i = 0; i < ropes.length; i++) {
            pq.add(ropes[i]);
        }
        int cost = 0;
        while (pq.size() > 1) {
            int min = pq.remove();
            int min2 = pq.remove();
            cost += min + min2;
            pq.add(min + min2);
        }
        return cost;
    }
}
