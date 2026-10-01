// Ported from the owner's Graphs/DisjointSetUnion.java (find with path compression, union by
// rank). The arrays belong to the instance and are sized n; the first link
// whose ends already share a root is the answer.
class Solution {
    private int[] par;
    private int[] rank;

    public int[] extraLink(int n, int[][] edges) {
        par = new int[n];
        rank = new int[n];
        for (int i = 0; i < n; i++) {
            par[i] = i;
        }
        for (int[] e : edges) {
            if (find(e[0]) == find(e[1])) {
                return e;
            }
            union(e[0], e[1]);
        }
        return new int[0];
    }

    private int find(int x) {
        if (x == par[x]) {
            return x;
        }
        return par[x] = find(par[x]);
    }

    private void union(int a, int b) {
        int parA = find(a);
        int parB = find(b);
        if (rank[parA] == rank[parB]) {
            par[parB] = parA;
            rank[parA]++;
        } else if (rank[parA] < rank[parB]) {
            par[parA] = parB;
        } else {
            par[parB] = parA;
        }
    }
}
