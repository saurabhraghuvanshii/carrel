import java.util.*;

// Ported from the owner's Graphs/TopologicalSorting.java (calcIndeg and topSort, the in-degree
// version). It collects the order it printed, and returns an empty order when
// it cannot take every course.
class Solution {
    static class Edge {
        int src;
        int dest;
        int wt;

        Edge(int s, int d, int w) {
            this.src = s;
            this.dest = d;
            this.wt = w;
        }
    }

    @SuppressWarnings("unchecked")
    private static ArrayList<Edge>[] build(int n, int[][] edges) {
        ArrayList<Edge>[] graph = new ArrayList[n];
        for (int i = 0; i < n; i++) {
            graph[i] = new ArrayList<>();
        }
        for (int[] e : edges) {
            graph[e[0]].add(new Edge(e[0], e[1], 1));
        }
        return graph;
    }

    public int[] order(int n, int[][] before) {
        ArrayList<Edge>[] graph = build(n, before);
        int[] indeg = new int[graph.length];
        for (int i = 0; i < graph.length; i++) {
            for (int j = 0; j < graph[i].size(); j++) {
                indeg[graph[i].get(j).dest]++;
            }
        }
        Queue<Integer> q = new LinkedList<>();
        for (int i = 0; i < indeg.length; i++) {
            if (indeg[i] == 0) {
                q.add(i);
            }
        }
        int[] order = new int[n];
        int taken = 0;
        while (!q.isEmpty()) {
            int curr = q.remove();
            order[taken++] = curr;
            for (int i = 0; i < graph[curr].size(); i++) {
                Edge e = graph[curr].get(i);
                indeg[e.dest]--;
                if (indeg[e.dest] == 0) {
                    q.add(e.dest);
                }
            }
        }
        return taken == n ? order : new int[0];
    }
}
