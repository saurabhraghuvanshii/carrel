#include <vector>
using namespace std;

static TreeNode* build(const vector<int>& nums, int lo, int hi) {
    if (lo > hi) return nullptr;
    int mid = lo + (hi - lo) / 2;
    TreeNode* n = new TreeNode(nums[mid]);
    n->left = build(nums, lo, mid - 1);
    n->right = build(nums, mid + 1, hi);
    return n;
}

TreeNode* buildTree(vector<int>& nums) {
    return build(nums, 0, (int)nums.size() - 1);
}
