#include <algorithm>
using namespace std;

static int depthAndBest(TreeNode* n, int& best) {
    if (!n) return 0;
    int l = depthAndBest(n->left, best), r = depthAndBest(n->right, best);
    best = max(best, l + r + 1);
    return max(l, r) + 1;
}

int longestPath(TreeNode* root) {
    int best = 0;
    depthAndBest(root, best);
    return best;
}
