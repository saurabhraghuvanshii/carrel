#include <algorithm>
using namespace std;

static int gain(TreeNode* n, int& best) {
    if (!n) return 0;
    int l = max(0, gain(n->left, best)), r = max(0, gain(n->right, best));
    best = max(best, n->val + l + r);
    return n->val + max(l, r);
}

int bestPathSum(TreeNode* root) {
    int best = root->val;
    gain(root, best);
    return best;
}
