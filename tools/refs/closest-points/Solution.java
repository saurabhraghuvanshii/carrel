import java.util.*;

// Ported from the owner's Heap/NearbyCar.java: the Point class and the body of
// main, now returning the positions of the k nearest instead of printing them.
class Solution {
    static class Point implements Comparable<Point> {
        int x;
        int y;
        int distSq;
        int idx;

        Point(int x, int y, int distSq, int idx) {
            this.x = x;
            this.y = y;
            this.distSq = distSq;
            this.idx = idx;
        }

        @Override
        public int compareTo(Point p2) {
            return this.distSq - p2.distSq;
        }
    }

    public int[] closest(int[][] pts, int k) {
        PriorityQueue<Point> pq = new PriorityQueue<>();
        for (int i = 0; i < pts.length; i++) {
            int distSq = pts[i][0] * pts[i][0] + pts[i][1] * pts[i][1];
            pq.add(new Point(pts[i][0], pts[i][1], distSq, i));
        }
        int[] nearest = new int[k];
        for (int i = 0; i < k; i++) {
            nearest[i] = pq.remove().idx;
        }
        return nearest;
    }
}
